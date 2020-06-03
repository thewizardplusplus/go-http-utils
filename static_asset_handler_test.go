package httputils

import (
	stdlog "log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-log/log"
	"github.com/go-log/log/print"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ExampleStaticAssetHandler() {
	// use the standard logger for error handling
	logger := stdlog.New(os.Stderr, "", stdlog.LstdFlags)
	staticAssetHandler := StaticAssetHandler(
		http.Dir("/var/www/example.com"),
		// wrap the standard logger via the github.com/go-log/log package
		print.New(logger),
	)

	http.Handle("/", staticAssetHandler)
	stdlog.Fatal(http.ListenAndServe(":8080", nil))
}

func TestStaticAssetHandler(test *testing.T) {
	type fileSystemComponents struct {
		fileInfos  []os.FileInfo
		files      []http.File
		fileSystem http.FileSystem
	}
	type args struct {
		fileSystemComponents fileSystemComponents
		logger               log.Logger
	}
	type handlerArgs struct {
		writer  http.ResponseWriter
		request *http.Request
	}

	for _, data := range []struct {
		name        string
		args        args
		handlerArgs handlerArgs
		wantHeader  http.Header
	}{
		{
			name: "success without redirecting to the index.html file",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					fileInfo := new(MockFileInfo)
					fileInfo.On("IsDir").Return(false)
					fileInfo.On("Name").Return("file")
					fileInfo.
						On("ModTime").
						Return(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))
					fileInfo.On("Size").Return(int64(4))

					file := new(MockFile)
					file.On("Stat").Return(fileInfo, nil)
					file.
						On("Read", mock.MatchedBy(func([]byte) bool { return true })).
						Return(func(buffer []byte) int { return copy(buffer, "test") }, nil)
					file.On("Close").Return(nil)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/path/to/file").Return(file, nil)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{fileInfo},
						files:      []http.File{file},
						fileSystem: fileSystem,
					}
				}(),
				logger: new(MockLogger),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusOK).Return()
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/path/to/file",
					nil,
				),
			},
			wantHeader: http.Header{
				"Accept-Ranges":  {"bytes"},
				"Content-Length": {"4"},
				"Content-Type":   {"text/plain"},
				"Last-Modified":  {"Mon, 02 Jan 2006 15:04:05 GMT"},
			},
		},
		{
			name: "success with redirecting to the index.html file",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					directoryInfo := new(MockFileInfo)
					directoryInfo.On("IsDir").Return(true)

					fileInfo := new(MockFileInfo)
					fileInfo.On("IsDir").Return(false)
					fileInfo.On("Name").Return("file")
					fileInfo.
						On("ModTime").
						Return(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))
					fileInfo.On("Size").Return(int64(4))

					directory := new(MockFile)
					directory.On("Stat").Return(directoryInfo, nil)
					directory.On("Close").Return(nil)

					file := new(MockFile)
					file.On("Stat").Return(fileInfo, nil)
					file.
						On("Read", mock.MatchedBy(func([]byte) bool { return true })).
						Return(func(buffer []byte) int { return copy(buffer, "test") }, nil)
					file.On("Close").Return(nil)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/").Return(directory, nil)
					fileSystem.On("Open", "/index.html").Return(file, nil)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{directoryInfo, fileInfo},
						files:      []http.File{directory, file},
						fileSystem: fileSystem,
					}
				}(),
				logger: new(MockLogger),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusOK).Return()
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
				request: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodGet,
						"http://example.com/path/to/file",
						nil,
					)
					request.Header.Set("Accept", "text/html")

					return request
				}(),
			},
			wantHeader: http.Header{
				"Accept-Ranges":  {"bytes"},
				"Content-Length": {"4"},
				"Content-Type":   {"text/plain"},
				"Last-Modified":  {"Mon, 02 Jan 2006 15:04:05 GMT"},
			},
		},
		{
			name: "error with the http.File.Stat() method",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					fileInfo := new(MockFileInfo)

					file := new(MockFile)
					file.On("Stat").Return(nil, iotest.ErrTimeout)
					file.On("Close").Return(nil)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/path/to/file").Return(file, nil)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{fileInfo},
						files:      []http.File{file},
						fileSystem: fileSystem,
					}
				}(),
				logger: new(MockLogger),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusInternalServerError).Return()
					writer.On("Write", []byte("500 Internal Server Error\n")).Return(26, nil)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/path/to/file",
					nil,
				),
			},
			wantHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			name: "error with the http.File.Read() method",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					fileInfo := new(MockFileInfo)
					fileInfo.On("IsDir").Return(false)
					fileInfo.On("Name").Return("file")
					fileInfo.
						On("ModTime").
						Return(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))
					fileInfo.On("Size").Return(int64(4))

					file := new(MockFile)
					file.On("Stat").Return(fileInfo, nil)
					file.
						On("Read", mock.MatchedBy(func([]byte) bool { return true })).
						Return(0, iotest.ErrTimeout)
					file.On("Close").Return(nil)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/path/to/file").Return(file, nil)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{fileInfo},
						files:      []http.File{file},
						fileSystem: fileSystem,
					}
				}(),
				logger: new(MockLogger),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusOK).Return()

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/path/to/file",
					nil,
				),
			},
			wantHeader: http.Header{
				"Accept-Ranges":  {"bytes"},
				"Content-Length": {"4"},
				"Content-Type":   {"text/plain"},
				"Last-Modified":  {"Mon, 02 Jan 2006 15:04:05 GMT"},
			},
		},
		{
			name: "error with the http.File.Close() method",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					fileInfo := new(MockFileInfo)
					fileInfo.On("IsDir").Return(false)
					fileInfo.On("Name").Return("file")
					fileInfo.
						On("ModTime").
						Return(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))
					fileInfo.On("Size").Return(int64(4))

					file := new(MockFile)
					file.On("Stat").Return(fileInfo, nil)
					file.
						On("Read", mock.MatchedBy(func([]byte) bool { return true })).
						Return(func(buffer []byte) int { return copy(buffer, "test") }, nil)
					file.On("Close").Return(iotest.ErrTimeout)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/path/to/file").Return(file, nil)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{fileInfo},
						files:      []http.File{file},
						fileSystem: fileSystem,
					}
				}(),
				logger: new(MockLogger),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusOK).Return()
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/path/to/file",
					nil,
				),
			},
			wantHeader: http.Header{
				"Accept-Ranges":  {"bytes"},
				"Content-Length": {"4"},
				"Content-Type":   {"text/plain"},
				"Last-Modified":  {"Mon, 02 Jan 2006 15:04:05 GMT"},
			},
		},
		{
			name: "error with the http.FileSystem.Open() method",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					fileInfo := new(MockFileInfo)

					file := new(MockFile)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/path/to/file").Return(nil, iotest.ErrTimeout)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{fileInfo},
						files:      []http.File{file},
						fileSystem: fileSystem,
					}
				}(),
				logger: new(MockLogger),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusInternalServerError).Return()
					writer.On("Write", []byte("500 Internal Server Error\n")).Return(26, nil)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/path/to/file",
					nil,
				),
			},
			wantHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			name: "error with the http.ResponseWriter.Write() method",
			args: args{
				fileSystemComponents: func() fileSystemComponents {
					fileInfo := new(MockFileInfo)
					fileInfo.On("IsDir").Return(false)
					fileInfo.On("Name").Return("file")
					fileInfo.
						On("ModTime").
						Return(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC))
					fileInfo.On("Size").Return(int64(4))

					file := new(MockFile)
					file.On("Stat").Return(fileInfo, nil)
					file.
						On("Read", mock.MatchedBy(func([]byte) bool { return true })).
						Return(func(buffer []byte) int { return copy(buffer, "test") }, nil)
					file.On("Close").Return(nil)

					fileSystem := new(MockFileSystem)
					fileSystem.On("Open", "/path/to/file").Return(file, nil)

					return fileSystemComponents{
						fileInfos:  []os.FileInfo{fileInfo},
						files:      []http.File{file},
						fileSystem: fileSystem,
					}
				}(),
				logger: func() log.Logger {
					logger := new(MockLogger)
					logger.
						On(
							"Logf",
							mock.MatchedBy(func(string) bool { return true }),
							iotest.ErrTimeout,
						).
						Return()

					return logger
				}(),
			},
			handlerArgs: handlerArgs{
				writer: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Header").Return(http.Header{"Content-Type": {"text/plain"}})
					writer.On("WriteHeader", http.StatusOK).Return()
					writer.On("Write", []byte("test")).Return(0, iotest.ErrTimeout)

					return writer
				}(),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/path/to/file",
					nil,
				),
			},
			wantHeader: http.Header{
				"Accept-Ranges":  {"bytes"},
				"Content-Length": {"4"},
				"Content-Type":   {"text/plain"},
				"Last-Modified":  {"Mon, 02 Jan 2006 15:04:05 GMT"},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := StaticAssetHandler(
				data.args.fileSystemComponents.fileSystem,
				data.args.logger,
			)
			handler.ServeHTTP(data.handlerArgs.writer, data.handlerArgs.request)

			for _, fileInfo := range data.args.fileSystemComponents.fileInfos {
				mock.AssertExpectationsForObjects(test, fileInfo)
			}
			for _, file := range data.args.fileSystemComponents.files {
				mock.AssertExpectationsForObjects(test, file)
			}
			mock.AssertExpectationsForObjects(
				test,
				data.args.fileSystemComponents.fileSystem,
				data.args.logger,
				data.handlerArgs.writer,
			)
			assert.Equal(test, data.wantHeader, data.handlerArgs.writer.Header())
		})
	}
}
