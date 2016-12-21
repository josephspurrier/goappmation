@ECHO OFF

SET GOVERSION=1.7.4
SET OLDPORTVERSION=1.7.1-r.1
SET NEWPORTVERSION=%GOVERSION%-r.1
SET LITEIDEVERSION=30.2

REM Download Go 32-bit
goappmation.exe -version=%GOVERSION% "Go-windows-386.json"

REM Download Go 64-bit
goappmation.exe -version=%GOVERSION% "Go-windows-amd64.json"

REM Download latest Go Portable
goappmation.exe -version=%OLDPORTVERSION% "Go-portable.json"

REM Download LiteIDE
goappmation.exe -version=%LITEIDEVERSION% "LiteIDE.json"

REM Rename the folder
MOVE "GoPortWin%OLDPORTVERSION%" "GoPortWin%NEWPORTVERSION%"

REM Remove the go directory
RMDIR /S /Q "GoPortWin%NEWPORTVERSION%\go"

REM Remove the liteide directory
RMDIR /S /Q "GoPortWin%NEWPORTVERSION%\liteide"

REM Copy 32-bit files
XCOPY /E /H /Q /Y "Go v%GOVERSION%-386\go" "GoPortWin%NEWPORTVERSION%\go\"

REM Copy 64-bit files
XCOPY /E /H /Q /Y "Go v%GOVERSION%-amd64\go" "GoPortWin%NEWPORTVERSION%\go\"

REM Copy LiteIDE
XCOPY /E /H /Q /Y "LiteIDEWinX%LITEIDEVERSION%" "GoPortWin%NEWPORTVERSION%\liteide\"

REM Remove the directories
RMDIR /S /Q "Go v%GOVERSION%-386"
RMDIR /S /Q "Go v%GOVERSION%-amd64"
RMDIR /S /Q "LiteIDEWinX%LITEIDEVERSION%"

PAUSE
