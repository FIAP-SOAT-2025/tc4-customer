# language: pt
Funcionalidade: Gerenciamento de Cliente
  Como um sistema de gerenciamento de clientes
  Eu quero validar os dados dos clientes
  Para garantir a integridade das informações armazenadas

  Cenário: Criar um cliente válido
    Dado que eu tenho os seguintes dados do cliente:
      | nome     | John Doe            |
      | cpf      | 111.444.777-35      |
      | email    | john@example.com    |
    Quando eu criar um novo cliente
    Então o cliente deve ser criado com sucesso
    E o cliente deve ter um ID gerado
    E o cliente deve ter data de criação preenchida
    E o cliente deve ter data de atualização preenchida

  Cenário: Tentar criar cliente com nome vazio
    Dado que eu tenho os seguintes dados do cliente:
      | nome     |                     |
      | cpf      | 11144477735         |
      | email    | john@example.com    |
    Quando eu criar um novo cliente
    Então deve retornar um erro
    E o código do erro deve ser "NAME_EMPTY"
    E o cliente não deve ser criado

  Cenário: Tentar criar cliente com CPF inválido
    Dado que eu tenho os seguintes dados do cliente:
      | nome     | John Doe            |
      | cpf      | 12345678901         |
      | email    | john@example.com    |
    Quando eu criar um novo cliente
    Então deve retornar um erro
    E o código do erro deve ser "INVALID_CPF"
    E o cliente não deve ser criado

  Cenário: Tentar criar cliente com email inválido
    Dado que eu tenho os seguintes dados do cliente:
      | nome     | John Doe            |
      | cpf      | 11144477735         |
      | email    | invalid-email       |
    Quando eu criar um novo cliente
    Então deve retornar um erro
    E o código do erro deve ser "INVALID_EMAIL"
    E o cliente não deve ser criado

  Cenário: Atualizar apenas o nome do cliente
    Dado que existe um cliente com os dados:
      | nome     | John Doe            |
      | cpf      | 11144477735         |
      | email    | john@example.com    |
    Quando eu atualizar o nome para "Jane Doe"
    Então a atualização deve ser bem-sucedida
    E o nome do cliente deve ser "Jane Doe"
    E a data de atualização deve ser posterior à data inicial

  Cenário: Atualizar apenas o email do cliente
    Dado que existe um cliente com os dados:
      | nome     | John Doe            |
      | cpf      | 11144477735         |
      | email    | john@example.com    |
    Quando eu atualizar o email para "jane@example.com"
    Então a atualização deve ser bem-sucedida
    E o email do cliente deve ser "jane@example.com"
    E a data de atualização deve ser posterior à data inicial

  Cenário: Atualizar nome e email do cliente
    Dado que existe um cliente com os dados:
      | nome     | John Doe            |
      | cpf      | 11144477735         |
      | email    | john@example.com    |
    Quando eu atualizar os dados:
      | nome     | Jane Smith              |
      | email    | jane.smith@example.com  |
    Então a atualização deve ser bem-sucedida
    E o nome do cliente deve ser "Jane Smith"
    E o email do cliente deve ser "jane.smith@example.com"
    E a data de atualização deve ser posterior à data inicial

  Cenário: Tentar atualizar cliente com nome vazio
    Dado que existe um cliente com os dados:
      | nome     | John Doe            |
      | cpf      | 11144477735         |
      | email    | john@example.com    |
    Quando eu atualizar o nome para ""
    Então deve retornar um erro
    E o nome do cliente não deve ser alterado

  Cenário: Tentar atualizar cliente com email inválido
    Dado que existe um cliente com os dados:
      | nome     | John Doe            |
      | cpf      | 11144477735         |
      | email    | john@example.com    |
    Quando eu atualizar o email para "invalid"
    Então deve retornar um erro
    E o email do cliente não deve ser alterado
