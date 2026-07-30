

# RadOTP
## Acerca del Proyecto
RadOTP es un servidor Radius dedicado, diseñado específicamente para servir como fuente de autenticación de dos factores para usuarios de VPN en firewalls como Fortigate y Cisco ASA.

```
  +----------------------------------+
  | internet                         |
  |                                  |
  |   +--------------------------+   |
  |   | User                     |   |
  |   |--------------------------|   |
  |   | Username: [YourUsername] |   |
  |   | Password: [YourPassword] |   |
  |   | OTP:      [YourOTP]      |   |
  |   +--------------------------+   |
  +----------------------------------+

                 |
                 v
 
  +----------------------------------+
  | Corporate Network                |
  |                                  |
  |   +--------------------------+   |
  |   | Fortigate SSL VPN        |   |
  |   |--------------------------|   |
  |   | Connects to RadOTP       |   |
  |   +--------------------------+   |
  |              |                   |
  |              v                   |
  |   +--------------------------+   |
  |   | RadOTP Radius Server     |   |
  |   |--------------------------|   |
  |   | 1. AD Password Check     |   |
  |   | 2. OTP Validation        |   |
  |   +--------------------------+   |
  +----------------------------------+

```

 **Cuenta con las siguientes características**:

-  Soporta autenticación de dos factores (2FA) utilizando contraseñas de un solo uso (OTP) que se generan y almacenan en una base de datos PostgreSQL.

-  Puede conectarse a Microsoft Active Directory mediante un cliente LDAP y verificar las credenciales de los usuarios.

-  Cuenta con una interfaz web que permite a los administradores gestionar usuarios, ver registros y configurar ajustes.

-  Exponen la actividad de los usuarios a Prometheus y Grafana para fines de monitoreo y alertas.

-  Utiliza Access-Challenge de Radius para interactuar con los usuarios y solicitar información adicional.

-  Ofrece alta disponibilidad y rendimiento mediante la replicación de PostgreSQL y el despliegue con docker-compose.

-  Dispone de una API REST que permite a aplicaciones externas gestionar usuarios de forma programática.


**El servidor Radius admite cuatro modos de autenticación:**
  
-  only_password: Este modo autentica a los usuarios contra un servidor LDAP/LDAPS de Active Directory. Los usuarios solo necesitan ingresar su contraseña de AD para iniciar sesión.

-  only_otp: Este modo autentica a los usuarios únicamente con una base de datos de OTP. Los usuarios solo necesitan ingresar un código de contraseña de un solo uso (OTP) para iniciar sesión.

-  two_fa: Este modo habilita la autenticación de dos factores (2FA). Los usuarios deben ingresar tanto su contraseña de AD como un código OTP para iniciar sesión.

-  two_fa_optional_otp: Este modo es similar a two_fa, pero solo aplica 2FA a los usuarios que tienen un OTP en la base de datos. Los usuarios que no tengan un OTP pueden iniciar sesión únicamente con su contraseña de AD.

![RadOTP](https://github.com/Abbas-gheydi/radotp/blob/main/assets/radotp.jpg)

### Cómo instalar
Primero, instale Docker y docker-compose, luego descargue [RadOTP](https://github.com/Abbas-gheydi/radotp/releases/latest) y despliéguelo mediante:  
```bash
docker-compose up -d  
```

### Cómo usarlo:
  
Para configurar el servidor Radius, siga estos pasos:

•  Edite el archivo config/radiusd.conf y realice los siguientes cambios:

•  Reemplace la dirección del servidor LDAP con la dirección de su controlador de dominio (servidor LDAP).

•  Reemplace el secreto de Radius con una frase de contraseña fuerte y segura.

•  Establezca Authentication_Mode en la sección de radius según sus necesidades.

•  Edite el archivo docker-compose.yml y realice los siguientes cambios:

•  Reemplace el nombre de usuario y la contraseña de la base de datos con nuevas credenciales.

•  Asegúrese de que el nombre de usuario y la contraseña de la base de datos coincidan con los del archivo config/radiusd.conf.

•  Guarde los cambios y reinicie los contenedores de Docker.   


• En su navegador, escriba http**s**://IP_ADDRESS y use la contraseña admin/admin para iniciar sesión y gestionar usuarios.   


• Descargue [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en_US&gl=US) y escanee el código QR.   


• Configure los firewalls de Cisco o Fortinet para utilizar radOTP (Radius) como fuente de autenticación:   
[Fortigate](https://docs.fortinet.com/document/fortigate/6.0.0/cookbook/200757/connecting-the-fortigate-to-the-radius-server)   
[Cisco ASA](https://www.cisco.com/c/en/us/support/docs/security/asa-5500-x-series-next-generation-firewalls/98594-configure-radius-authentication.html)   


## Licencia

MIT
