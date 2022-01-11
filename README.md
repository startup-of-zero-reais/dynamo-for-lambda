# Dynamo for lambda

Pacote criado para utilizar o dynamo db com lambda functions

# Para rodar o projeto

Certifique-se de ter instalado Docker e docker-compose.
Em seguida rode o seguinte comando:

```shell
docker-compose up -d
```

Isso fará com que os containers do Dynamo Local e documentação
auto gerada do Godoc subam

Após isso você terá disponível:
- DynamoDB - [http://localhost:8000](http://localhost:8000)
- Godoc - [http://localhost:6060](http://localhost:6060/pkg/github.com/startup-of-zero-reais/dynamo-for-lambda)

# Testes

Rode os testes e veja quanto coverage tem

```shell
./tests
```