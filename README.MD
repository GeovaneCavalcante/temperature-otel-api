
# Temperature API - README.md

## Introdução

Bem-vindo ao repositório da Temperature API, uma solução simples e eficiente para gerenciar e acessar dados de temperatura. Este README serve como um guia para configurar e utilizar a aplicação.

## Pré-requisitos

Antes de começar, certifique-se de ter as seguintes ferramentas instaladas:

- Docker
- Make
- Git (para clonar o repositório)

## Configuração

Siga os passos abaixo para configurar o ambiente do projeto.

### Clonando o Repositório

Primeiro, clone o repositório da Temperature API para a sua máquina local:

```bash
git clone https://github.com/GeovaneCavalcante/temperature-otel-api
cd temperature-otel-api
```


### Subindo o Projeto

Para subir o projeto, você tem duas opções:


1. **Construa a imagem e suba o projeto:**

   Siga os comandos abaixo para construir a imagem e subir o projeto:

   ```bash
   make docker-build-image
   make docker-up
   ```



## Executando Testes

Para executar os testes da aplicação, utilize o comando:

```bash
make test
```

Para gerar a cobertura de código e um relatório em `coverage.html`, execute:

```bash
make test-coverage
```

Após a execução, o arquivo `coverage.html` será gerado na raiz do projeto, permitindo que você visualize detalhadamente a cobertura de testes.

## Uso da API

A Temperature API está disponível no seguinte link: http://localhost:8080/temperature?cep=65909001

Para utilizar, passe o CEP desejado via query params, por exemplo: `http://localhost:8080/temperature?cep=SEU_CEP_AQUI`.

Para acessar os traces da aplicação acesse: http://localhost:9411/zipkin/
