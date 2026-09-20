# Clima por CEP com Observabilidade

Projeto desenvolvido como desafio prático do Labs Go, módulo do curso Go Expert da Full Cycle.

Sistema distribuído em Go para consulta de clima por CEP com dois serviços, OpenTelemetry, OTEL Collector e Zipkin.

## Arquitetura

- **Serviço A:** recebe `POST /`, valida o CEP e chama o Serviço B.
- **Serviço B:** recebe `GET /clima-por-cep/{cep}`, consulta ViaCEP e WeatherAPI e realiza as conversões.
- **OTEL Collector:** recebe traces via OTLP gRPC e os envia ao Zipkin.
- **Zipkin:** visualiza o fluxo distribuído.

O Serviço A é o único serviço exposto para as requisições da aplicação. Internamente, ele chama o Serviço B.

## Configuração

Copie o arquivo de exemplo e informe uma chave do [WeatherAPI](https://www.weatherapi.com/):

```bash
cp .env.example .env
```

Edite `.env`:

```env
WEATHER_API_KEY=sua-chave-do-weatherapi
```

## Execução

Suba os quatro componentes:

```bash
docker compose up --build
```

Envie uma requisição para o Serviço A:

```bash
curl -X POST http://localhost:8080 \
  -H 'Content-Type: application/json' \
  -d '{"cep":"01001000"}'
```

Resposta esperada:

```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

As conversões seguem as fórmulas do desafio:

- Fahrenheit: `C × 1.8 + 32`;
- Kelvin: `C + 273`.

## Erros

- `422 invalid zipcode`: CEP que não é uma string com exatamente oito dígitos;
- `404 can not find zipcode`: CEP válido, mas não encontrado.

## Observabilidade

Acesse o Zipkin em [http://localhost:9411](http://localhost:9411) e procure o fluxo completo:

`Request → Serviço A → Serviço B`

Também devem aparecer os spans manuais das consultas ao ViaCEP e ao WeatherAPI.

## Testes

```bash
go test ./...
```
