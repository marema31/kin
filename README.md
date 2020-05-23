# kin
Automatic home page for docker hosted web applications.

Serve a templated web site with variables populated from specific labels on containers. The list of variables are refresh every 10 seconds.
Default site delivered by `kin`can be replaced by your own site.

`kin` is perfect for development environment or home-lab to quickly provides a dynamic bookmarks to all containers availables on the docker daemon. It is not meant to be use in production environment. 

`kin` will not proxy your site, but you can use it alongside [Traefik](https://containo.us/traefik/) to have fully automatic configuration.

## Usage

`kin [Flags]`

### Configuration
The configuration can be done through flags, environment variables or entry in a configuration file (in JSON, YAML or TOML format).

Flags (short form)             | Env var      |File entry | meaning
-------------------------------|--------------|-----------|----------------------------------------------------------
  `--base <URL path> (-b)`     |`KIN_BASE`    |`base`     | Base URL (default "/")
  `--config <path> (-c)  `     |`KIN_CONFIG`  |           | Config file without extension (default is $HOME/.kin.yaml)
  `--debug (-d) / --quiet (-q)`|`KIN_LOGLEVEL`|`log.level`| Log more information / Log only errors
  `--json (-j) `               |`KIN_JSON`    |`log.json` | In present logs are JSON-formatted
  `--logpath <path> (-l) `     |`KIN_LOGPATH` |`log.path` | Log file path (default "-" for screen)
  `--port int (-p)       `     |`KIN_PORT`    |`port`     | Port to listen (default 8080)
  `--root <path> (-r)    `     |`KIN_ROOT`    |`root`     | Template root path (default is $HOME/.kin_root)
  `--swarm               `     |`KIN_SWARM`   |`swarm`    | Docker swarm


### Site structure
If template root path is provided to kin, this can contains as much file neededs by your site. At each request, `kin` lookup for the corresponding file with `.tpl` extension, render the template and serve it. If a template is not found, it will serve the file without any rendering.

Template are not limited to HTML file, it could be any text format (Javascript, JSON, CSV, etc)

See the `site` folder in this repository for an example of site structure.

### Template files
For template, `kin` provides :
    * `.Containers` : array of list of containers information indexed by `kin_group` label value. Only containers with `kin_name` label will be added to the list.
    * `.Environments`: array of environment variables indexed by the environment variable name.

Each entry of `.Containers` is an list of all the container informations with the same `kin_group` label.
    
Container informations correspond to the labels found on the containers with `kin_name` label. To use a attribute, the template should only contains `{{.AttributeName}}`. 

 Available attributes for a container are: 

Attribute | Label
----------|-----------
`Group`   | `kin_name`
`Name`    | `kin_group`
`Type`    | `kin_type`
`URL`     | `kin_url`

All containers with `kin_name` label and without `kin_group` label will be placed in `""` entry of `.Containers`.

Kin use Golang templating system that provides [more features](https://golang.org/pkg/html/template/).

#### Example using groups
```html
 <body>
     <h1> {{index .Environments "USER"}}'s home lab</h1>
     <ul>
        {{ range $group, $containers := .Containers }}
        <li> {{ $group }}
          <ul>
            {{ range $containers}}
            <li><a href="{{.URL}}" class="{{.Type}}">{{.Name}}</li>
            {{end}}
          </ul>
        </>
        {{end}}
    </ul>
  </body>
```

#### Example without groups
```html
 <body>
     <ul>
         {{ range index .Containers "" }}
         <li><a href="{{.URL}}" class="{{.Type}}">{{.Name}}</li>
        {{end}}
    </ul>
  </body>
```
#### Adding labels on container: docker run
```bash
docker run --rm --label 'kin_name=front' --label 'kin_group=dev' --label 'kin_type=web' --label 'kin_url=http://localhost:1234/my-front'  nginx
```

#### Adding labels on container: docker-compose
```yaml
front: 
  image: nginx
  labels:
    - "kin_name=front"
    - "kin_group=dev"
    - "kin_type=web"
    - "kin_url=http://localhost:1234/my-front"
```
## Installation

### Via Homebrew

You can install `kin` with [Homebrew](https://github.com/marema31/homebrew-marema31):

```bash
brew tap marema31/marema31
brew install kin
```

### Via Docker
There is a [Docker image](https://hub.docker.com/r/marema31/kin/) that you can use to run `kin` in a container:
```bash
docker pull marema31/kin
docker run -it --rm -p 8080:8080 -v "/var/run/docker.sock:/var/run/docker.sock:ro" marema31/kin
```

### Via prebuilt binaries
You can download prebuilt binaries from the [Release page](https://github.com/marema31/kin/releases)

### From source

If you want to build `kin` from source, you need Golang 1.14 or
higher. You can build everything:

```bash
make
```

## Contribution
I've made this project as a real use case to learn Golang.
I've tried to adopt the Go mindset but I'm sure that other gophers could do better. 

If this project can be useful for you, feel free to open issues, propose documentation updates or even pull request.

Contributions are of course always welcome!

1. Fork marema31/kin (https://github.com/marema31/kin/fork)
2. Create a feature branch
3. Commit your changes
4. Create a Pull Request

See [`CONTRIBUTING.md`](https://github.com/marema31/kin/blob/master/CONTRIBUTING.md) for details.

## License

Copyright (c) 2020-present [Marc Carmier](https://github.com/marema31)

Licensed under [BSD-2 Clause License](./LICENSE)