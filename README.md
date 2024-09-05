### Resolución:
```bash
./generar-compose.sh docker-compose-dev.yaml 5
```
Para correr el script, y para levantar la app:
`make docker-compose-up`
Se hizo un script que genera el docker compose a través de unos archivos de texto y un script de python que se encarga de generar el docker-compose.yml.

### Ejercicio N°1:
Además, definir un script de bash `generar-compose.sh` que permita crear una definición de DockerCompose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```