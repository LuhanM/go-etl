## GO-ETL

Esta aplicação recebe, tranforma e persiste registro através de uma operação REST.

# Como usar 

## Docker

1 - Clone o repositório para a sua máquina;
2 - Acesso o diretório criado pelo clone;
3 - Execute o comando "docker-compose up -d --build"
    obs: Executando este comando o script "create-table.sql" será executado automáticamente, porém,
        caso esteja utilizando outra instação do postgres será necessário executá-lo e também alterar as variáveis de ambiente no compose. 
4 - Para importar o arquivo, através do Postman ou aplicação semelhante execute o verbo POST para localhost:8080/arquivo tendo, no form-data as seguintes chaves:
     _________________________________________________________
    |      KEY:type       |          VALUE                    |
    |---------------------|------------------------------------
    |      file: Tile     | o arquivo que deseja importar     | 
    |---------------------|-----------------------------------|
    | transacional: Texto | true or false: define se será     |
    |                     | utilizada ou não transação global |
    |_____________________|___________________________________|

5 - Caso importação ocorra com sucesso será retorndo status 200 e os registrões estarão no banco e caso contrário 400.


