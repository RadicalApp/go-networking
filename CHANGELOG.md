# Change Log
All notable changes to this project will be documented in this file.

`go-networking` closely adheres to [Semantic Versioning](http://semver.org/).

--- 

## [0.2.1](https://github.com/RadicalApp/go-networking/releases/tag/0.2.1) / 05th'April 2017
* Expose `FileWrapper` struct
* Cleaned up code and added comments

## [0.2.0](https://github.com/RadicalApp/go-networking/releases/tag/0.2.0) / 21st'March 2017

#### Changed
* `Response` is a complex `struct`.
	* Contains `data`, original http `request` & http `response`.
* `ConnectionState` are now `string` type.


## [0.1.1](https://github.com/RadicalApp/go-networking/releases/tag/0.1.1) / 21st'March 2017

#### Changed
* All RESTful requests run asynchronously.

#### Updated
* Updated test cases.
	* Can test async calls with [Goblin](https://github.com/franela/goblin).

#### Removed
* Removed `Completion` type. 
	* Pass a function as a parameter instead to ensure no name clashes.


## [0.1.0](https://github.com/RadicalApp/go-networking/releases/tag/0.1.0) / 20th'March 2017

* Initial implementation with `GET`, `POST` and `UPLOAD` method


<!-- 
SAMPLE CHANGE LOG!
#### Added
* 
 * 


#### Updated
* 
 * 


#### Changed
* 
 * 


#### Fixed
* 
 * 


#### Removed
* 
 *  
-->
