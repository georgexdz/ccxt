
import os
import sys
import re
import humps

import esprima


EX_NAME = 'Kucoin'


def get_ex_list():
    return [
        'kucoin'
    ]


def get_func_code_map(code_str):
    return []


def init_code():
    return '''
package {ex}

import (
	. "ccxt-master/go/base"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type {ex.capitalize()} struct {
	Exchange
}

    '''


JS_FOLD = './js'
def read_code_str(ex):
    with open(os.path.join(JS_FOLD, '{ex}.js')) as f:
        f.read


FUNC_ARG_MAP = {
    'createOrder': 'symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}',
    'fetchBalance': 'params map[string]interface{}',
    'cancelOrder': 'id string, symbol string, params map[string]interface{}',
    'fetchOpenOrders': 'symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchOrdersByStatus': 'status string, symbol string, since int64, limit int64, params map[string]interface{}',
    'fetchOrderBook': 'symbol string, limit int64, params map[string]interface{}',
    'fetchOrder': 'id string, symbol string, params map[string]interface{}',
    'sign': 'id string, symbol string, params map[string]interface{}',
    'parseOrder': 'id string, symbol string, params map[string]interface{}',
}
RETURN_MAP = {
    'createOrder': 'order map[string]interface{}, err error',
    'fetchBalance': 'balanceResult *Account, err error',
    'cancelOrder': 'response interface{}, err error',
    'fetchOrdersByStatus': 'orders []*Order, err error',
    'fetchOpenOrders': 'orders []*Order, err error',
    'fetchOrderBook': 'orderBook *OrderBook, err error',
    'fetchOrder': 'order *Order, err error',
    'sign': 'order *Order, err error',
    'parseOrder': 'order interface{}',
}
NIL_MAP = {
    'string': '""',
    'int64': '0',
    'error': 'nil'
}
NO_ERROR_RETURN_FUNCS = ('safe', 'tostring', 'milliseconds', 'uuid', 'account', 'ifthenelse', 'iso8601')
SIDE = None


def get_return(func_name):
    return RETURN_MAP.get(func_name, None)


def get_arg(func_name):
    default = '(interface{})'
    return FUNC_ARG_MAP.get(func_name, default)


DEFAULT_FUNC_ARGS = {
    'SafeValue'.lower(): (3, 'nil'),
    'SafeString'.lower(): (3, '""'),
    'ApiFunc'.lower(): (4, 'nil'),
    'SafeFloat'.lower(): (3, 0),
    'SafeInteger'.lower(): (3, 0),
}


def ThisExpression(syntax, info={}):
    return 'self'


def MemberExpression(syntax, info={}):
    method_name = call_func_by_syntax(syntax.property)

    info['error_check'] = False
    if syntax.object.type == 'ThisExpression' and syntax.property.type == 'Identifier' and method_need_check_err(method_name):
        info['error_check'] = True

    if syntax.object.type == 'ThisExpression':
        method_name = humps.pascalize(method_name)

    obj = call_func_by_syntax(syntax.object)
    if syntax.object.type == 'ThisExpression':
        return f'{obj}.{method_name}'
    else:
        m = {
            'toString': f'fmt.Sprintf("%v", {obj})',
            'length': f'self.Length({obj})',
        }
        default = f'{obj}[{method_name}]'

        if syntax.property.type == 'Identifier' and syntax.property.name == 'split':
            print('xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')
            print(info)
            print(syntax)
            if 'arg_str' in info:
                return f'strings.Split({obj}, {info["arg_str"]})'

        if info.get('pre') in ['right', 'init', 'test']:
            default = f'self.Member({obj}, {method_name})'
        if info.get('pre') in ['left']:
            default = f'self.SetValue({obj}, {method_name})'

        return m.get(syntax.property.name, default)


def CallExpression(syntax, info={}):
    arg_str = ','.join([call_func_by_syntax(arg, info) for arg in syntax.arguments])
    # default arg
    if syntax.callee.property.name.lower() in DEFAULT_FUNC_ARGS:
        arg_info = DEFAULT_FUNC_ARGS[syntax.callee.property.name.lower()]
        arg_str += f', {arg_info[1]}' * (arg_info[0] - len(syntax.arguments))
    info.update({'arg_str': arg_str})

    pre_part = f'{call_func_by_syntax(syntax.callee, info)}'

    # toString -> fmt.Sprintf
    if '(' in pre_part and ')' in pre_part:
        return pre_part

    call_str = call_func_by_syntax(syntax.callee, info)
    if re.findall(r'self\.(Private|Public)', call_str):
        info = DEFAULT_FUNC_ARGS['apifunc']
        arg_str += f', {info[1]}' * (info[0] - len(syntax.arguments) - 1)
        return f'self.ApiFunc("{syntax.callee.property.name}", {arg_str})'
    # xx.split('_') -> string.Split(xx, '_')
    elif not call_str.endswith(')'):
        return f'{call_str}({arg_str})'
    else:
        return call_str


def ExpressionStatement(syntax, info={}):
    if syntax.expression.type == 'CallExpression':
        return f'_, err = {call_func_by_syntax(syntax.expression)}\n if err != nil {{\n return \n}}'
    else:
        return call_func_by_syntax(syntax.expression)


def method_need_check_err(name, info={}):
    if any(name.lower().startswith(o) for o in NO_ERROR_RETURN_FUNCS):
        return False
    return True


def VariableDeclarator(syntax, info={}):
    left = call_func_by_syntax(syntax.id)
    operator = ':='
    if left in ARG_TYPE or left in RET_TYPE:
        operator = '='

    info = {'pre': 'init'}
    if syntax.init.type == 'Identifier' and syntax.init.name == 'undefined':
        return f'var {left} interface{{}}'
    if syntax.id.type == 'ArrayPattern':
        return f'{left} {operator} self.Unpack{len(syntax.id.elements)}({call_func_by_syntax(syntax.init, info)})'

    right = call_func_by_syntax(syntax.init, info)

    if info.get('error_check'):
        return f'{left}, err {operator} {right}\nif err!= nil {{\nreturn\n}}'
    else:
        return f'{left} {operator} {right}'


def VariableDeclaration(syntax, info={}):
    return '\n'.join([call_func_by_syntax(o) for o in syntax.declarations])


def ObjectExpression(syntax, info={}):
    k_v = ''
    for one in syntax.properties:
        k_v += f'\n{call_func_by_syntax(one.key)}: {call_func_by_syntax(one.value)},'
    ret = f'map[string]interface{{}}{{{k_v}\n}}'
    return ret


def Literal(syntax, info={}):
    val = syntax.raw.replace("'", '"')
    return val


def Identifier(syntax, info={}):
    m = {
        'undefined': 'nil',
        'type': 'typ',
    }
    return m.get(syntax.name, syntax.name)


def BinaryExpression(syntax, info={}):
    operator = syntax.operator

    if syntax.right.name == 'undefined':
        if operator == '===':
            return f'self.TestNil({call_func_by_syntax(syntax.left)})'
        if operator == '!==':
            return f'!self.TestNil({call_func_by_syntax(syntax.left)})'

    if operator == '===':
        operator = '=='
    if operator == '!==':
        operator = '!='
    if operator == 'in':
        # return f'_, ok := {call_func_by_syntax(syntax.right)}[{call_func_by_syntax(syntax.left)}]; ok'
        return f'self.InMap({call_func_by_syntax(syntax.left)}, {call_func_by_syntax(syntax.right)})'
    left = call_func_by_syntax(syntax.left)
    right = call_func_by_syntax(syntax.right)

    if left in ARG_TYPE and right == 'nil':
        right = NIL_MAP.get(ARG_TYPE[left], 'nil')

    return f'{left} {operator} {right}'


def call_func_by_syntax(syntax, info={}):
    return globals()[syntax.type](syntax, info)


def AssignmentExpression(syntax, info={}):
    operator = syntax.operator
    if syntax.left.type == 'MemberExpression':
        arg1 = call_func_by_syntax(syntax.left.object)
        arg2 = call_func_by_syntax(syntax.left.property)
        arg3 = call_func_by_syntax(syntax.right, {"pre": "right"})
        return f'self.SetValue({arg1}, {arg2}, {arg3})'
    else:
        return f'{call_func_by_syntax(syntax.left, {"pre": "left"})} {operator} {call_func_by_syntax(syntax.right, {"pre": "right"})}'


def BlockStatement(syntax, info={}):
    lines = '\n'.join([call_func_by_syntax(block) for block in syntax.body])
    return f'{{\n{lines}\n}}'


def IfStatement(syntax, info={}):
    ret = f'if {call_func_by_syntax(syntax.test)} '

    ret = f'if self.ToBool({call_func_by_syntax(syntax.test)}) '
    ret += f'{call_func_by_syntax(syntax.consequent)}'

    if syntax.alternate:
        ret += (
            f' else '
            f'{call_func_by_syntax(syntax.alternate)}'
        )

    return ret


def ConditionalExpression(syntax, info={}):
    return f'self.IfThenElse(self.ToBool({call_func_by_syntax(syntax.test, {"pre": "test"})}), {call_func_by_syntax(syntax.consequent)}, {call_func_by_syntax(syntax.alternate)})'


def ArrayPattern(syntax, info={}):
    return ', '.join([call_func_by_syntax(e) for e in syntax.elements])


def UnaryExpression(syntax, info={}):
    return f'{syntax.operator}self.ToBool({call_func_by_syntax(syntax.argument)})'


def ReturnStatement(syntax, info={}):
    if info.get('return') and ',' in info['return']:
        return f'return {call_func_by_syntax(syntax.argument)}, nil'
    else:
        return f'return {call_func_by_syntax(syntax.argument)}'


def UpdateExpression(syntax, info={}):
    return f'{call_func_by_syntax(syntax.argument)}{syntax.operator}'


def ForStatement(syntax, info={}):
    return f'for {call_func_by_syntax(syntax.init)}; {call_func_by_syntax(syntax.test)}; {call_func_by_syntax(syntax.update)} {call_func_by_syntax(syntax.body)}'


def ThrowStatement(syntax, info={}):
    err_info = ','.join([call_func_by_syntax(o) for o in syntax.argument.arguments])
    return f'err = errors.New({err_info});\nreturn'


def LogicalExpression(syntax, info={}):
    return f'{call_func_by_syntax(syntax.left)} {syntax.operator} {call_func_by_syntax(syntax.right)}'


def ArrayExpression(syntax, info={}):
    elements = ','.join([call_func_by_syntax(o) for o in syntax.elements])
    return f'[]interface{{}}{{{elements}}}'


def capital_first(s, info={}):
    return s[0].upper() + s[1:]


FUNC_LINES = 99
ARG_TYPE = dict()
RET_TYPE = dict()
def FunctionDeclaration(syntax, info={}):
    result_code = ''
    func_name = syntax.id.name

    global ARG_TYPE, RET_TYPE

    func_arg_str = get_arg(func_name)
    a = [tuple(pair.split(' ')) for pair in func_arg_str.split(', ')]
    ARG_TYPE = dict(a)
    func_ret_str = get_return(func_name)
    a = [tuple(pair.split(' ')) for pair in func_ret_str.split(', ')]
    RET_TYPE = dict(a)

    if syntax.body.type == 'BlockStatement':
        for idx, block in enumerate(syntax.body.body):
            result_code += '\n' + call_func_by_syntax(block, {'return': func_ret_str}) + '\n'
            if idx == FUNC_LINES:
                break

    return f'''func (self *{EX_NAME}) {capital_first(func_name)} ({func_arg_str}) ({func_ret_str}) {{
    {result_code}
}}
    '''


def syntax_analysis(syntax, info={}):
    ret = ''
    for idx, func in enumerate(syntax):
        ret += call_func_by_syntax(func)
    return ret

def parse_by_syntax(str_code):
    str_code = str_code.replace('await ', '')
    str_code = str_code.replace('async ', 'function ')
    syntax = esprima.parse(str_code)
    return syntax_analysis(syntax.body)


def test():
    global FUNC_LINES
    FUNC_LINES = 99
    code = '''
    async createOrder (symbol, type, side, amount, price = undefined, params = {}) {
        await this.loadMarkets ();
        const marketId = this.marketId (symbol);
        // required param, cannot be used twice
        const clientOid = this.uuid ();
        const request = {
            'clientOid': clientOid,
            'side': side,
            'symbol': marketId,
            'type': type,
        };
        if (type !== 'market') {
            request['price'] = this.priceToPrecision (symbol, price);
            request['size'] = this.amountToPrecision (symbol, amount);
        } else {
            if (this.safeValue (params, 'quoteAmount')) {
                // used to create market order by quote amount - https://github.com/ccxt/ccxt/issues/4876
                request['funds'] = this.amountToPrecision (symbol, amount);
            } else {
                request['size'] = this.amountToPrecision (symbol, amount);
            }
        }
        const response = await this.privatePostOrders (this.extend (request, params));
        //
        //     {
        //         code: '200000',
        //         data: {
        //             "orderId": "5bd6e9286d99522a52e458de"
        //         }
        //    }
        //
        const data = this.safeValue (response, 'data', {});
        const timestamp = this.milliseconds ();
        const order = {
            'id': this.safeString (data, 'orderId'),
            'symbol': symbol,
            'type': type,
            'side': side,
            'price': price,
            'cost': undefined,
            'filled': undefined,
            'remaining': undefined,
            'timestamp': timestamp,
            'datetime': this.iso8601 (timestamp),
            'fee': undefined,
            'status': 'open',
            'clientOrderId': clientOid,
            'info': data,
        };
        if (!this.safeValue (params, 'quoteAmount')) {
            order['amount'] = amount;
        }
        return order;
    }
    '''
    print(parse_by_syntax(code))

    FUNC_LINES = 8
    code = '''
    async fetchBalance (params = {}) {
        await this.loadMarkets ();
        let type = undefined;
        const request = {};
        if ('type' in params) {
            type = params['type'];
            if (type !== undefined) {
                request['type'] = type;
            }
            params = this.omit (params, 'type');
        } else {
            const options = this.safeValue (this.options, 'fetchBalance', {});
            type = this.safeString (options, 'type', 'trade');
        }
        const response = await this.privateGetAccounts (this.extend (request, params));
        //
        //     {
        //         "code":"200000",
        //         "data":[
        //             {"balance":"0.00009788","available":"0.00009788","holds":"0","currency":"BTC","id":"5c6a4fd399a1d81c4f9cc4d0","type":"trade"},
        //             {"balance":"3.41060034","available":"3.41060034","holds":"0","currency":"SOUL","id":"5c6a4d5d99a1d8182d37046d","type":"trade"},
        //             {"balance":"0.01562641","available":"0.01562641","holds":"0","currency":"NEO","id":"5c6a4f1199a1d8165a99edb1","type":"trade"},
        //         ]
        //     }
        //
        const data = this.safeValue (response, 'data', []);
        const result = { 'info': response };
        for (let i = 0; i < data.length; i++) {
            const balance = data[i];
            const balanceType = this.safeString (balance, 'type');
            if (balanceType === type) {
                const currencyId = this.safeString (balance, 'currency');
                const code = this.safeCurrencyCode (currencyId);
                const account = this.account ();
                account['total'] = this.safeFloat (balance, 'balance');
                account['free'] = this.safeFloat (balance, 'available');
                account['used'] = this.safeFloat (balance, 'holds');
                result[code] = account;
            }
        }
        return this.parseBalance (result);
    }
'''
    print(parse_by_syntax(code))

    code = '''
        async fetchOrder (id, symbol = undefined, params = {}) {
        await this.loadMarkets ();
        const request = {
            'orderId': id,
        };
        let market = undefined;
        if (symbol !== undefined) {
            market = this.market (symbol);
        }
        const response = await this.privateGetOrdersOrderId (this.extend (request, params));
        const responseData = response['data'];
        return this.parseOrder (responseData, market);
    }
'''
    FUNC_LINES = 99
    print(parse_by_syntax(code))
    code = '''
        async cancelOrder (id, symbol = undefined, params = {}) {
        const request = { 'orderId': id };
        const response = await this.privateDeleteOrdersOrderId (this.extend (request, params));
        return response;
    }
'''
    FUNC_LINES = 99
    print(parse_by_syntax(code))

    code = '''
        async fetchOrderBook (symbol, limit = undefined, params = {}) {
        const level = this.safeInteger (params, 'level', 2);
        let levelLimit = level.toString ();
        if (levelLimit === '2') {
            if (limit !== undefined) {
                if ((limit !== 20) && (limit !== 100)) {
                    throw new ExchangeError (this.id + ' fetchOrderBook limit argument must be undefined, 20 or 100');
                }
                levelLimit += '_' + limit.toString ();
            }
        }
        await this.loadMarkets ();
        const marketId = this.marketId (symbol);
        const request = { 'symbol': marketId, 'level': levelLimit };
        const response = await this.publicGetMarketOrderbookLevelLevel (this.extend (request, params));
        //
        // 'market/orderbook/level2'
        // 'market/orderbook/level2_20'
        // 'market/orderbook/level2_100'
        //
        //     {
        //         "code":"200000",
        //         "data":{
        //             "sequence":"1583235112106",
        //             "asks":[
        //                 // ...
        //                 ["0.023197","12.5067468"],
        //                 ["0.023194","1.8"],
        //                 ["0.023191","8.1069672"]
        //             ],
        //             "bids":[
        //                 ["0.02319","1.6000002"],
        //                 ["0.023189","2.2842325"],
        //             ],
        //             "time":1586584067274
        //         }
        //     }
        //
        // 'market/orderbook/level3'
        //
        //     {
        //         "code":"200000",
        //         "data":{
        //             "sequence":"1583731857120",
        //             "asks":[
        //                 // id, price, size, timestamp in nanoseconds
        //                 ["5e915f8acd26670009675300","6925.7","0.2","1586585482194286069"],
        //                 ["5e915f8ace35a200090bba48","6925.7","0.001","1586585482229569826"],
        //                 ["5e915f8a8857740009ca7d33","6926","0.00001819","1586585482149148621"],
        //             ],
        //             "bids":[
        //                 ["5e915f8acca406000ac88194","6925.6","0.05","1586585482384384842"],
        //                 ["5e915f93cd26670009676075","6925.6","0.08","1586585491334914600"],
        //                 ["5e915f906aa6e200099b49f6","6925.4","0.2","1586585488941126340"],
        //             ],
        //             "time":1586585492487
        //         }
        //     }
        //
        const data = this.safeValue (response, 'data', {});
        const timestamp = this.safeInteger (data, 'time');
        const orderbook = this.parseOrderBook (data, timestamp, 'bids', 'asks', level - 2, level - 1);
        orderbook['nonce'] = this.safeInteger (data, 'sequence');
        return orderbook;
    }'''
    FUNC_LINES = 99
    print(parse_by_syntax(code))

    code = '''
        async fetchOrdersByStatus (status, symbol = undefined, since = undefined, limit = undefined, params = {}) {
        await this.loadMarkets ();
        const request = {
            'status': status,
        };
        let market = undefined;
        if (symbol !== undefined) {
            market = this.market (symbol);
            request['symbol'] = market['id'];
        }
        if (since !== undefined) {
            request['startAt'] = since;
        }
        if (limit !== undefined) {
            request['pageSize'] = limit;
        }
        const response = await this.privateGetOrders (this.extend (request, params));
        //
        //     {
        //         code: '200000',
        //         data: {
        //             "currentPage": 1,
        //             "pageSize": 1,
        //             "totalNum": 153408,
        //             "totalPage": 153408,
        //             "items": [
        //                 {
        //                     "id": "5c35c02703aa673ceec2a168",   //orderid
        //                     "symbol": "BTC-USDT",   //symbol
        //                     "opType": "DEAL",      // operation type,deal is pending order,cancel is cancel order
        //                     "type": "limit",       // order type,e.g. limit,markrt,stop_limit.
        //                     "side": "buy",         // transaction direction,include buy and sell
        //                     "price": "10",         // order price
        //                     "size": "2",           // order quantity
        //                     "funds": "0",          // order funds
        //                     "dealFunds": "0.166",  // deal funds
        //                     "dealSize": "2",       // deal quantity
        //                     "fee": "0",            // fee
        //                     "feeCurrency": "USDT", // charge fee currency
        //                     "stp": "",             // self trade prevention,include CN,CO,DC,CB
        //                     "stop": "",            // stop type
        //                     "stopTriggered": false,  // stop order is triggered
        //                     "stopPrice": "0",      // stop price
        //                     "timeInForce": "GTC",  // time InForce,include GTC,GTT,IOC,FOK
        //                     "postOnly": false,     // postOnly
        //                     "hidden": false,       // hidden order
        //                     "iceberg": false,      // iceberg order
        //                     "visibleSize": "0",    // display quantity for iceberg order
        //                     "cancelAfter": 0,      // cancel orders time，requires timeInForce to be GTT
        //                     "channel": "IOS",      // order source
        //                     "clientOid": "",       // user-entered order unique mark
        //                     "remark": "",          // remark
        //                     "tags": "",            // tag order source
        //                     "isActive": false,     // status before unfilled or uncancelled
        //                     "cancelExist": false,   // order cancellation transaction record
        //                     "createdAt": 1547026471000  // time
        //                 },
        //             ]
        //         }
        //    }
        const responseData = this.safeValue (response, 'data', {});
        const orders = this.safeValue (responseData, 'items', []);
        return this.parseOrders (orders, market, since, limit);
    }
    '''
    print(parse_by_syntax(code))

    code = '''
        async fetchOpenOrders (symbol = undefined, since = undefined, limit = undefined, params = {}) {
        return await this.fetchOrdersByStatus ('active', symbol, since, limit, params);
    }
    '''
    print(parse_by_syntax(code))

    code = '''
    function parseOrder (order, market = undefined) {
        //
        // fetchOpenOrders, fetchClosedOrders
        //
        //     {
        //         "id": "5c35c02703aa673ceec2a168",   //orderid
        //         "symbol": "BTC-USDT",   //symbol
        //         "opType": "DEAL",      // operation type,deal is pending order,cancel is cancel order
        //         "type": "limit",       // order type,e.g. limit,markrt,stop_limit.
        //         "side": "buy",         // transaction direction,include buy and sell
        //         "price": "10",         // order price
        //         "size": "2",           // order quantity
        //         "funds": "0",          // order funds
        //         "dealFunds": "0.166",  // deal funds
        //         "dealSize": "2",       // deal quantity
        //         "fee": "0",            // fee
        //         "feeCurrency": "USDT", // charge fee currency
        //         "stp": "",             // self trade prevention,include CN,CO,DC,CB
        //         "stop": "",            // stop type
        //         "stopTriggered": false,  // stop order is triggered
        //         "stopPrice": "0",      // stop price
        //         "timeInForce": "GTC",  // time InForce,include GTC,GTT,IOC,FOK
        //         "postOnly": false,     // postOnly
        //         "hidden": false,       // hidden order
        //         "iceberg": false,      // iceberg order
        //         "visibleSize": "0",    // display quantity for iceberg order
        //         "cancelAfter": 0,      // cancel orders time，requires timeInForce to be GTT
        //         "channel": "IOS",      // order source
        //         "clientOid": "",       // user-entered order unique mark
        //         "remark": "",          // remark
        //         "tags": "",            // tag order source
        //         "isActive": false,     // status before unfilled or uncancelled
        //         "cancelExist": false,   // order cancellation transaction record
        //         "createdAt": 1547026471000  // time
        //     }
        //
        let symbol = undefined;
        const marketId = this.safeString (order, 'symbol');
        if (marketId !== undefined) {
            if (marketId in this.markets_by_id) {
                market = this.markets_by_id[marketId];
                symbol = market['symbol'];
            } else {
                const [ baseId, quoteId ] = marketId.split ('-');
                const base = this.safeCurrencyCode (baseId);
                const quote = this.safeCurrencyCode (quoteId);
                symbol = base + '/' + quote;
            }
            market = this.safeValue (this.markets_by_id, marketId);
        }
        if (symbol === undefined) {
            if (market !== undefined) {
                symbol = market['symbol'];
            }
        }
        const orderId = this.safeString (order, 'id');
        const type = this.safeString (order, 'type');
        const timestamp = this.safeInteger (order, 'createdAt');
        const datetime = this.iso8601 (timestamp);
        let price = this.safeFloat (order, 'price');
        const side = this.safeString (order, 'side');
        const feeCurrencyId = this.safeString (order, 'feeCurrency');
        const feeCurrency = this.safeCurrencyCode (feeCurrencyId);
        const feeCost = this.safeFloat (order, 'fee');
        const amount = this.safeFloat (order, 'size');
        const filled = this.safeFloat (order, 'dealSize');
        const cost = this.safeFloat (order, 'dealFunds');
        const remaining = amount - filled;
        // bool
        let status = order['isActive'] ? 'open' : 'closed';
        status = order['cancelExist'] ? 'canceled' : status;
        const fee = {
            'currency': feeCurrency,
            'cost': feeCost,
        };
        if (type === 'market') {
            if (price === 0.0) {
                if ((cost !== undefined) && (filled !== undefined)) {
                    if ((cost > 0) && (filled > 0)) {
                        price = cost / filled;
                    }
                }
            }
        }
        const clientOrderId = this.safeString (order, 'clientOid');
        return {
            'id': orderId,
            'clientOrderId': clientOrderId,
            'symbol': symbol,
            'type': type,
            'side': side,
            'amount': amount,
            'price': price,
            'cost': cost,
            'filled': filled,
            'remaining': remaining,
            'timestamp': timestamp,
            'datetime': datetime,
            'fee': fee,
            'status': status,
            'info': order,
            'lastTradeTimestamp': undefined,
            'average': undefined,
            'trades': undefined,
        };
    }'''
    print(parse_by_syntax(code))

    code = '''
    function sign (path, api = 'public', method = 'GET', params = {}, headers = undefined, body = undefined) {
        //
        // the v2 URL is https://openapi-v2.kucoin.com/api/v1/endpoint
        //                                †                 ↑
        //
        const versions = this.safeValue (this.options, 'versions', {});
        const apiVersions = this.safeValue (versions, api);
        const methodVersions = this.safeValue (apiVersions, method, {});
        const defaultVersion = this.safeString (methodVersions, path, this.options['version']);
        const version = this.safeString (params, 'version', defaultVersion);
        params = this.omit (params, 'version');
        let endpoint = '/api/' + version + '/' + this.implodeParams (path, params);
        const query = this.omit (params, this.extractParams (path));
        let endpart = '';
        headers = (headers !== undefined) ? headers : {};
        if (Object.keys (query).length) {
            if (method !== 'GET') {
                body = this.json (query);
                endpart = body;
                headers['Content-Type'] = 'application/json';
            } else {
                endpoint += '?' + this.urlencode (query);
            }
        }
        const url = this.urls['api'][api] + endpoint;
        if (api === 'private') {
            this.checkRequiredCredentials ();
            const timestamp = this.nonce ().toString ();
            headers = this.extend ({
                'KC-API-KEY': this.apiKey,
                'KC-API-TIMESTAMP': timestamp,
                'KC-API-PASSPHRASE': this.password,
            }, headers);
            const payload = timestamp + method + endpoint + endpart;
            const signature = this.hmac (this.encode (payload), this.encode (this.secret), 'sha256', 'base64');
            headers['KC-API-SIGN'] = this.decode (signature);
        }
        return { 'url': url, 'method': method, 'body': body, 'headers': headers };
    }
    '''
    FUNC_LINES= 8
    print(parse_by_syntax(code))

def translate():
    for ex in get_ex_list():
        str_code = read_code_str(ex)


if __name__ == '__main__':
    test()
