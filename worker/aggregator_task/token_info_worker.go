package aggregator_task

// TODO:
// token 信息定时维护补全.
// 通过 token 表值不存在来判定如
// SELECT *
// FROM token
// WHERE token_name = ''
//    OR token_symbol = ''
//    OR token_decimal = ''
//    OR token_logo = '';
// 为 token 生成 guid并写入，然后补全所有信息。
