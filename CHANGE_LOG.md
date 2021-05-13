# Change Log

## [v1.2](https://github.com/thewizardplusplus/go-http-utils/tree/v1.2) (2021-05-13)

- analogs:
  - adding an analog of the `http.Error()` function with the additional improvements:
    - additional logging of the error;
    - accepting of an error object instead of an error string;
- adding the function to extract the parameter with the specified name from the path part of the request URL and then scan it into the data;
- JSON:
  - adding the function to read bytes from the reader and then unmarshal them into the data;
  - adding the function to marshal the data and then write it in the writer:
    - additional setting of the corresponding content type;
    - additional setting of the specified status code.

## [v1.1.1](https://github.com/thewizardplusplus/go-http-utils/tree/v1.1.1) (2021-02-05)

- adding a simplified interface of the `http.Client` structure for mocking purposes.

## [v1.1](https://github.com/thewizardplusplus/go-http-utils/tree/v1.1) (2020-06-03)

- analogs:
  - adding an analog of the `http.FileServer()` function with applied `SPAFallbackMiddleware()` and `CatchingMiddleware()` middlewares.

## [v1.0](https://github.com/thewizardplusplus/go-http-utils/tree/v1.0) (2020-05-14)
