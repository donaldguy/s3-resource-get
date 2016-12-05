# s3-resource-get

A simple, interactively usable version of [s3-resource](https://github.com/s3-resource)'s [in](https://github.com/concourse/s3-resource/blob/master/cmd/in/main.go) using a .flyrc target

## Installing / Building

`go get github.com/donaldguy/flight-tracker`

or download from [Releases](/releases)

## Usage

`s3-resource-get` borrows its auth from the local `~/.flyrc`, so before you do anything
`fly login` if needed

```
s3-resource-get -t local pipeline/resource
```

or

```
s3-resource-get -t local -o ~/myresource.txt pipeline/resource
```

## Args

#### `-t, --target=` **(required)**
The same things you would pass to `fly`. You can use `fly login` (or copy the resulting `~/.flyrc`) to make sure you are authed

#### `-o, --output=` *(optional)*
A local path to save the file. If not specified will write to current working directory with name matching the basename of the remote file in s3

## License
Copyright (C) 2016 Donald Guy, Tulip Interfaces

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this project except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
