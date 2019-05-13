# stock

实时查看股票价格

## 终端
```vim
stock 股票代码（多个股票使用,号分割）

eg:stock  sh600271,sz300481
```

## vim插件

```vim
function! Stockprice(symbol)
  execute system('stock -f vim ' . a:symbol)
endfunction

map <leader>htxx :call Stockprice('sh600271')<cr>
```
