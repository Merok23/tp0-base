### Resolución:
`make docker-compose-up` Para levantar la app

Se hizo uso de de threading.lock con el context manager de python para proteger la sección crítica de la persistencia de los mensajes, el contador de agencias terminadas, y el levantado de los datos de las agencias.

Python protege sus estados internos con el GIL, sin embargo, el código que maneja variables como la cuenta que lleva el "contador de agencias terminadas" y el "diccionario de agencias" no está protegido por el GIL, y por eso se usa el lock `https://stackoverflow.com/questions/40072873/why-do-we-need-locks-for-threads-if-we-have-gil`.

Tambien se crearon levanto un thread por cliente, y se agrego todo lo relacionado a mantener esos threads (una funcion de reap y la lista de threads) para poder liberar los recursos de los threads una vez que terminan.

Opino que hubiese sido mejor usar `asyncio` para tareas y manejar un threadpool en vez de un thread por coneccion pero no se llego por tiempo.

## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

