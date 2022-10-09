//
// 	Copyright 2017 by marmot author: gdccmcm14@live.com.
// 	Licensed under the Apache License, Version 2.0 (the "License");
// 	you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
// 		http://www.apache.org/licenses/LICENSE-2.0
// 	Unless required by applicable law or agreed to in writing, software
// 	distributed under the License is distributed on an "AS IS" BASIS,
// 	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// 	See the License for the specific language governing permissions and
// 	limitations under the License
//

/*
Package miner is the core of this project, use to request for http api.

Example:

	package main

	import (
	"fmt"

	"github.com/admpub/marmot/miner"
	)

	func main() {
		// Use Default Worker, You can Also New One:
		// worker:=miner.New(nil)
		miner.SetLogLevel(miner.DEBUG)
		_, err := miner.SetURL("https://www.whitehouse.gov").Go()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(miner.String())
		}
	}
*/
package miner
