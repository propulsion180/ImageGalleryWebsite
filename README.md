# ImageGalleryWebsite

## What is the purpose of this project

To create a website in which I could practice using the HTMX library and learn how to use cookies to keep someone logged in.

## What is in this website

The main page contains buttons which will send you to a admin page, login page, and a sign up page as well as all the images in the database.

When you click on an image it will take you to a details page which contains the image and all the settings used to capture it along with a description and a location

The login and sign up page have obvious functions.

The admin page lets you add images (through a file chooser), set the settings you used to capture it (ISO, shutter speed, and aperture), set the location, and write a description of the image.

## Server Endpoints

### `GET /`

Returns the index.html when requested.

### `GET /all`

Returns a list of all of the imagedata objects stored in the database.

### `POST /image`

Returns the data of a single image based on the filepath given in the json body

```
{
    Filepath: string 
}
```

### `POST /addimage`

Takes in an form data which has all the fields for an ImageData. It also contains the image itself.

#### Form Data

| Field | Descripton |
| ----- | ---------- |
| `description` | Some text that describes the image |
| `location` | The location the image was taken, just a string |
| `iso` | The ISO the photo was taken at |
| `shutterSpeed` | The shutter speed the photo was taken at |
| `aperture` | The aperture the photo was taken at |
| `file` | The File that was selected |


#### Returns

| Reason | Code |
| ------ | ---- |
| Failed to parse form data, authorize the user  | 400 |
| Failed to get file from data, create the image directory, decode the image, create the outputfile, compress to webp, or fail to get the sub from jwt token | 500 |
| Failed to get the permissions of user from jwt token, not authorized | 401 |
| All good | 200 |

### `DELETE /delimage`

Send the filepath you want to be deleted then you will get a result.

```
{
    filepath: string
}
```

#### Returns

| Reason | Code |
| ------ | ---- |
| Failed to decode json data, get cookie, verify jwt token | 400 |
| Failed to get sub from jwt token, to delete entry from database | 500 |
| Failed to get perms from jwt claims | 401 |

### `PUT /setimage`

#### Form data

| Field | Descripton |
| ----- | ---------- |
| `description` | Some text that describes the image |
| `location` | The location the image was taken, just a string |
| `iso` | The ISO the photo was taken at |
| `shutterSpeed` | The shutter speed the photo was taken at |
| `aperture` | The aperture the photo was taken at |
| `filepath` | The Filepath of the image you want to modify the properties of. |

#### Returns 



| Reason | Code |
| ------ | ---- |
| filed to parse form data, failed to get cookie, failed to verify jwt token | 400 |
| failed to parse sub from jwt claims, to update the database | 500 |
| failed to get permisions from jwt token | 401 |
| Sucess | 200 |

### `POST /login`

```
{
    Username: string,
    Password: string,
}
```

#### Returns 

Adds http cookie to response.

```
{
    Message: string,
    Username: string,
    Admin: bool
}
```

| Reason | Code |
| ------ | ---- |
| Failed to decode json data | 400 |
| Failed to verify password with database, Failed to generate jwt token, failed to set jwt token | 500 |
| Invalid password | 401 |
| Sucess | 200 |


### `POST /tknlgn`

Used to restore state from a jwt token cookie.

#### Returns 

```
{
    Message: string,
    Userame: string,
    Admin: bool
}
```

| Reason | Code |
| ------ | ---- |
| Failed to get cookie from the request, Failed to verify jwt token, Failed to get perms from jwt token | 401 |
| Failed o get the sub from claims | 500 |
| Sucess | 200 |

### `POST /signup`

Takes same data as login as new credentials.
As you can see there is no way to add an admin without manually manipulating the database.

```
{
    Username: string,
    Password: string
}
```

#### Returns

| Reason | Code |
| ------ | ---- |
| Invalid payload | 400 | 
| Failed to add user | 401 | 
| Failed to add user | 500 | 
| Sucess | 200 |

### `POST /logout`

Takes the cookie from the request and invalidates it, removes the token from the database and returns the cookie.

This could just be a GET

#### Returns

| Reason | Code |
| ------ | ---- |
| Failed to verify jwt toke, failed to get claims, failed to delete token from database | 500 |
| Sucess | 200 | 




