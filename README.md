# bud
Example web service in Go using the seed, stem, and fertilize projects

# Example output

```
$ curl -X POST  --data '{"username": "user@example.com", "password": "gophers"}' http://localhost:8080/v1/UserService.Authenticate
{"token":"eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJidWQgcHJvamVjdCIsInN1YiI6ImYyZjFjZTM1LWM4OTAtNDM4ZS1iZDIzLWIyZmFlZDk5MGIzZSIsImV4cCI6MTY4NjQzMTkzMywiaWF0IjoxNjg2NDI4MzMzLCJyb2xlcyI6WyJBRE1JTiJdfQ.LdIGDioGtVLK1YRPMBMFqGdzJwd77MgTgCbhirl_rgWR33F17cX-7vDdY_B-bX1_gy-SeWnaLfK1zeoI2QTa8eAKVlQrf1Plx2l--o11xhDY880zOjO5SFy3MZ2exV5pQn3B6rz43M0qfn4u7MJHIvZEb9DObPRRq8PSCOpXRNFwAL6SOzpXHgawkBTmCMg6eN7cKcNrxygoHl14MFafdLt44Y5P-a7ly5QeHzeHqFm9g_rcGPrSgxRCP2IjmNj8inL6XxT53u2XFf1dypRxjNiLIh3gXJESEAkmtn90kIBc_dUDviL8r53yX0GfKSQnHn4efsXkERHxTtFStAifgw"}
```