# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Результат оптимизации
```
$ go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
File: ___86go_build_github_com_kuznet1_urlshrt_cmd_shortener
Build ID: 79e75763f78896b953f307653a99751b4302d000
Type: inuse_space
Time: 2025-11-06 02:19:56 +04
Duration: 60.05s, Total samples = 7181.51kB
Showing nodes accounting for -6161.27kB, 85.79% of 7181.51kB total
flat  flat%   sum%        cum   cum%
-3078kB 42.86% 42.86%    -3078kB 42.86%  runtime.allocm
-1024.64kB 14.27% 57.13% -1032.55kB 14.38%  github.com/jackc/pgx/v5.connect
1024.23kB 14.26% 42.87%  1024.23kB 14.26%  internal/profile.(*Profile).postDecode
-516.01kB  7.19% 50.05%  -516.01kB  7.19%  github.com/jackc/pgx/v5/internal/iobufpool.init.0.func1
-516.01kB  7.19% 57.24%  -516.01kB  7.19%  io.init.func1
-514kB  7.16% 64.39%     -514kB  7.16%  bufio.NewReaderSize (inline)
-512.69kB  7.14% 71.53%  -512.69kB  7.14%  vendor/golang.org/x/sys/cpu.initOptions
512.44kB  7.14% 64.40%   512.44kB  7.14%  crypto/tls.Client (inline)
-512.38kB  7.13% 71.53%  -512.38kB  7.13%  bytes.growSlice
-512.16kB  7.13% 78.66%  -512.16kB  7.13%  net/http.(*Request).WithContext (inline)
-512.05kB  7.13% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5/pgproto3.(*RowDescription).Decode
-512.03kB  7.13% 92.92%  -512.03kB  7.13%  syscall.anyToSockaddr
-512.02kB  7.13% 100.05%  -512.02kB  7.13%  github.com/jackc/pgx/v5.(*Conn).BeginTx
512.02kB  7.13% 92.92%   512.02kB  7.13%  github.com/jackc/pgx/v5/pgtype.NewMap
512.02kB  7.13% 85.79%   512.02kB  7.13%  net.sockaddrToTCP
0     0% 85.79%  -516.01kB  7.19%  bufio.(*Writer).Flush
0     0% 85.79%     -514kB  7.16%  bufio.NewReader (inline)
0     0% 85.79%  -512.38kB  7.13%  bytes.(*Buffer).Write
0     0% 85.79%  -512.38kB  7.13%  bytes.(*Buffer).grow
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).Handshake (inline)
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).HandshakeContext (inline)
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).Write
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).clientHandshake
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).handshakeContext
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).readHandshake
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).readHandshakeBytes
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).readRecord (inline)
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*Conn).readRecordOrCCS
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*clientHandshakeStateTLS13).handshake
0     0% 85.79%  -512.38kB  7.13%  crypto/tls.(*clientHandshakeStateTLS13).readServerCertificate
0     0% 85.79%  -512.05kB  7.13%  database/sql.(*Conn).QueryContext
0     0% 85.79%  -512.05kB  7.13%  database/sql.(*Conn).QueryRowContext (inline)
0     0% 85.79%   513.16kB  7.15%  database/sql.(*DB).BeginTx
0     0% 85.79%   513.16kB  7.15%  database/sql.(*DB).BeginTx.func1
0     0% 85.79% -2057.74kB 28.65%  database/sql.(*DB).QueryContext
0     0% 85.79% -2057.74kB 28.65%  database/sql.(*DB).QueryContext.func1
0     0% 85.79% -2057.74kB 28.65%  database/sql.(*DB).QueryRowContext (inline)
0     0% 85.79%   513.16kB  7.15%  database/sql.(*DB).begin
0     0% 85.79%  -512.02kB  7.13%  database/sql.(*DB).beginDC
0     0% 85.79%  -512.02kB  7.13%  database/sql.(*DB).beginDC.func1
0     0% 85.79% -1032.55kB 14.38%  database/sql.(*DB).conn
0     0% 85.79% -2057.74kB 28.65%  database/sql.(*DB).query
0     0% 85.79%  -512.05kB  7.13%  database/sql.(*DB).queryDC
0     0% 85.79%  -512.05kB  7.13%  database/sql.(*DB).queryDC.func1
0     0% 85.79% -1544.57kB 21.51%  database/sql.(*DB).retry
0     0% 85.79%  -512.02kB  7.13%  database/sql.ctxDriverBegin
0     0% 85.79%  -512.05kB  7.13%  database/sql.ctxDriverQuery
0     0% 85.79% -1024.07kB 14.26%  database/sql.withLock
0     0% 85.79% -2056.73kB 28.64%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
0     0% 85.79%   513.16kB  7.15%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
0     0% 85.79%  -512.05kB  7.13%  github.com/golang-migrate/migrate/v4.(*Migrate).Up
0     0% 85.79%  -512.05kB  7.13%  github.com/golang-migrate/migrate/v4/database/postgres.(*Postgres).Version
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5.(*Conn).Prepare
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5.(*Conn).Query
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5.(*Conn).getStatementDescription
0     0% 85.79% -1032.55kB 14.38%  github.com/jackc/pgx/v5.ConnectConfig
0     0% 85.79%  -516.01kB  7.19%  github.com/jackc/pgx/v5/internal/iobufpool.Get
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5/pgconn.(*PgConn).Prepare
0     0% 85.79%  -512.38kB  7.13%  github.com/jackc/pgx/v5/pgconn.(*PgConn).flushWithPotentialWriteReadDeadlock
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5/pgconn.(*PgConn).receiveMessage
0     0% 85.79%  -519.93kB  7.24%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
0     0% 85.79% -1032.02kB 14.37%  github.com/jackc/pgx/v5/pgconn.ParseConfigWithOptions.func1
0     0% 85.79%  -519.93kB  7.24%  github.com/jackc/pgx/v5/pgconn.connectOne
0     0% 85.79%  -519.93kB  7.24%  github.com/jackc/pgx/v5/pgconn.connectPreferred
0     0% 85.79%   512.44kB  7.14%  github.com/jackc/pgx/v5/pgconn.startTLS
0     0% 85.79%   516.01kB  7.19%  github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).bgRead
0     0% 85.79%  -512.38kB  7.13%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Flush
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive
0     0% 85.79% -1032.02kB 14.37%  github.com/jackc/pgx/v5/pgproto3.NewFrontend
0     0% 85.79% -1032.02kB 14.37%  github.com/jackc/pgx/v5/pgproto3.newChunkReader (inline)
0     0% 85.79%  -512.02kB  7.13%  github.com/jackc/pgx/v5/stdlib.(*Conn).BeginTx
0     0% 85.79%  -512.05kB  7.13%  github.com/jackc/pgx/v5/stdlib.(*Conn).QueryContext
0     0% 85.79% -1032.55kB 14.38%  github.com/jackc/pgx/v5/stdlib.(*driverConnector).Connect
0     0% 85.79%   513.16kB  7.15%  github.com/kuznet1/urlshrt/internal/handler.Handler.ShortenJSON
0     0% 85.79% -1544.57kB 21.51%  github.com/kuznet1/urlshrt/internal/middleware.(*Auth).Authentication-fm.(*Auth).Authentication.func1
0     0% 85.79% -1544.57kB 21.51%  github.com/kuznet1/urlshrt/internal/middleware.Compression.func1
0     0% 85.79% -1544.57kB 21.51%  github.com/kuznet1/urlshrt/internal/middleware.RequestLogger.Logging-fm.RequestLogger.Logging.func1
0     0% 85.79% -2057.74kB 28.65%  github.com/kuznet1/urlshrt/internal/repository.(*DBRepo).CreateUser
0     0% 85.79%   513.16kB  7.15%  github.com/kuznet1/urlshrt/internal/repository.(*DBRepo).Put
0     0% 85.79%  -512.05kB  7.13%  github.com/kuznet1/urlshrt/internal/repository.NewDBRepo
0     0% 85.79%  -512.05kB  7.13%  github.com/kuznet1/urlshrt/internal/repository.NewRepo
0     0% 85.79%  -512.05kB  7.13%  github.com/kuznet1/urlshrt/internal/repository.applyMigrations
0     0% 85.79%   513.16kB  7.15%  github.com/kuznet1/urlshrt/internal/service.(*Service).Shorten
0     0% 85.79%  -512.03kB  7.13%  internal/poll.(*FD).Accept
0     0% 85.79%  -512.03kB  7.13%  internal/poll.accept
0     0% 85.79%  1024.23kB 14.26%  internal/profile.Parse
0     0% 85.79%  1024.23kB 14.26%  internal/profile.parseUncompressed
0     0% 85.79%  -516.01kB  7.19%  io.Copy (inline)
0     0% 85.79%  -516.01kB  7.19%  io.CopyN
0     0% 85.79%  -516.01kB  7.19%  io.copyBuffer
0     0% 85.79%  -516.01kB  7.19%  io.discard.ReadFrom
0     0% 85.79% -1024.08kB 14.26%  main.main
0     0% 85.79%   512.02kB  7.13%  net.(*Dialer).DialContext
0     0% 85.79%  -512.03kB  7.13%  net.(*TCPListener).Accept
0     0% 85.79%  -512.03kB  7.13%  net.(*TCPListener).accept
0     0% 85.79%  -512.03kB  7.13%  net.(*netFD).accept
0     0% 85.79%   512.02kB  7.13%  net.(*netFD).dial
0     0% 85.79%   512.02kB  7.13%  net.(*sysDialer).dialParallel
0     0% 85.79%   512.02kB  7.13%  net.(*sysDialer).dialSerial
0     0% 85.79%   512.02kB  7.13%  net.(*sysDialer).dialSingle
0     0% 85.79%   512.02kB  7.13%  net.(*sysDialer).dialTCP
0     0% 85.79%   512.02kB  7.13%  net.(*sysDialer).doDialTCP (inline)
0     0% 85.79%   512.02kB  7.13%  net.(*sysDialer).doDialTCPProto
0     0% 85.79%   512.02kB  7.13%  net.internetSocket
0     0% 85.79%   512.02kB  7.13%  net.socket
0     0% 85.79%  1024.23kB 14.26%  net/http.(*ServeMux).ServeHTTP
0     0% 85.79%  -512.03kB  7.13%  net/http.(*Server).ListenAndServe
0     0% 85.79%  -512.03kB  7.13%  net/http.(*Server).Serve
0     0% 85.79%  -516.01kB  7.19%  net/http.(*chunkWriter).Write
0     0% 85.79%  -516.01kB  7.19%  net/http.(*chunkWriter).writeHeader
0     0% 85.79% -2062.51kB 28.72%  net/http.(*conn).serve
0     0% 85.79%  -516.01kB  7.19%  net/http.(*response).finishRequest
0     0% 85.79%  -520.34kB  7.25%  net/http.HandlerFunc.ServeHTTP
0     0% 85.79%  -512.03kB  7.13%  net/http.ListenAndServe (inline)
0     0% 85.79%     -514kB  7.16%  net/http.newBufioReader
0     0% 85.79% -1032.50kB 14.38%  net/http.serverHandler.ServeHTTP
0     0% 85.79%  1024.23kB 14.26%  net/http/pprof.Index
0     0% 85.79%  1024.23kB 14.26%  net/http/pprof.collectProfile
0     0% 85.79%  1024.23kB 14.26%  net/http/pprof.handler.ServeHTTP
0     0% 85.79%  1024.23kB 14.26%  net/http/pprof.handler.serveDeltaProfile
0     0% 85.79%  -512.69kB  7.14%  runtime.doInit (inline)
0     0% 85.79%  -512.69kB  7.14%  runtime.doInit1
0     0% 85.79%    -1026kB 14.29%  runtime.findRunnable
0     0% 85.79%     -513kB  7.14%  runtime.handoffp
0     0% 85.79%    -1026kB 14.29%  runtime.injectglist
0     0% 85.79%    -1026kB 14.29%  runtime.injectglist.func1
0     0% 85.79% -1536.77kB 21.40%  runtime.main
0     0% 85.79%    -3078kB 42.86%  runtime.mstart
0     0% 85.79%    -3078kB 42.86%  runtime.mstart0
0     0% 85.79%    -3078kB 42.86%  runtime.mstart1
0     0% 85.79%    -3078kB 42.86%  runtime.newm
0     0% 85.79%    -1539kB 21.43%  runtime.resetspinning
0     0% 85.79%     -513kB  7.14%  runtime.retake
0     0% 85.79%    -2565kB 35.72%  runtime.schedule
0     0% 85.79%    -3078kB 42.86%  runtime.startm
0     0% 85.79%     -513kB  7.14%  runtime.sysmon
0     0% 85.79%    -1539kB 21.43%  runtime.wakep
0     0% 85.79% -1032.02kB 14.37%  sync.(*Pool).Get
0     0% 85.79%  -512.03kB  7.13%  syscall.Accept4
0     0% 85.79%  -512.69kB  7.14%  vendor/golang.org/x/sys/cpu.init.0
```