# Token
## Summary
A library to generate and manage time-base access tokens.  Also works as a
simple password authentication database.
## Quickstart
```
store, _ := NewTokenStore("", "tokens", time.Second*1)
store.AddUser("test", "password")
token, err := store.NewToken("test", "password")
if store.IsValidToken(token) {
  fmt.Println("Valid!")
}
time.Sleep(time.Second*2)
if !store.IsValidToken(token) {
  fmt.Println("Invalid!")
}
```
