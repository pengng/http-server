package main

const directory = `
<!DOCTYPE html>
<html lang="zh-cn">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document</title>
  <style>
    h1 {
      overflow: auto;
      white-space: nowrap;
    }
    .table-container {
      overflow: auto;
    }
    table {
      width: 100%;
      text-align: left;
      border-collapse: collapse;
    }
    table th, table td {
      white-space: nowrap;
    }
    table tbody td {
      font-size: 14px;
      padding-right: 20px;
    }
    table thead tr {
      border-bottom: 2px solid #efefef;
    }
    table tbody tr:hover {
      background: #f5f5f5;
    }
    table tr {
      height: 40px;
      border-bottom: 1px solid #efefef;
    }
  </style>
</head>

<body>
  <h1>Index of {{.Pathname}}</h1>
  <div class="table-container">
    <table class="table table-striped table-hover">
      <thead>
        <tr>
          <th>修改日期</th>
          <th>大小</th>
          <th>文件名</th>
        </tr>
      </thead>
      <tbody>
        {{if ne .Pathname "/"}}
        <tr>
          <td></td>
          <td></td>
          <td>
            <a href="../">..</a>
          </td>
        </tr>
        {{end}}
        {{range .List}}
        <tr>
          <td>{{.ModTime}}</td>
          <td>{{.Size}}</td>
          <td>
            <a href="{{.Name}}">{{.Name}}</a>
          </td>
        </tr>
        {{end}}
      </tbody>
    </table>
  </div>
</body>

</html>
`