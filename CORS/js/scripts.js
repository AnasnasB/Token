
var Token = "";
async function requestget(url, data) {
	let res = await fetch(url, {
		method: 'POST', // *GET, POST, PUT, DELETE, etc.
		mode: 'cors', // no-cors, *cors, same-origin
		cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
		credentials: 'same-origin', // include, *same-origin, omit
		headers: {
		  'Content-Type': 'application/json',
		  // 'Content-Type': 'application/x-www-form-urlencoded',
		},
		redirect: 'follow', // manual, *follow, error
		referrerPolicy: 'no-referrer', // no-referrer, *client
		body: data // body data type must match "Content-Type" header
	  });
	let inf = await res;
	console.log(inf);
	
	let jso = await res.text();
	document.getElementById("resp").innerHTML = jso;
}

async function requestpost(url, data) {
	let res = await fetch( url, {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		credentials: 'same-origin',
		headers: {
			'Content-Type': 'application/json'
		},
		redirect: 'follow',
		referrerPolicy: 'no-referrer',
		body: data
	});
	let inf = await res;
	console.log(inf);
	
	let jso = await res.text();
	Token = jso;
	if (Token == "No acess") {
		Token = ""
		document.getElementById("resp").innerHTML = jso;
	}	
		
}
	


function getget() {
	var data2 = {
		token : Token 
	}
	var data3 = JSON.stringify(data2);
	requestget("http://localhost:8080/data" , data3);
}

function getpost() {
     var data1 = {
		 token: "",
		 name:""  ,
		 pass:""
		  };
	data1.name = String (document.getElementById("bodyusername").value)
	data1.pass = String (document.getElementById("bodypass").value)
	data1.token = Token
	var data = JSON.stringify(data1);
	
	requestpost("http://localhost:8080/authorization", data);
}
