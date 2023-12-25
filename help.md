``` 
docker build --tag translator:latest .
```


```
docker run --name translator -p 8987:8987 -d translator:latest
```