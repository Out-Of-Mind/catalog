// main section
const delete_btn = document.querySelectorAll(".main__group-item--btn");
const admin_delete_btn = document.querySelectorAll(".main__groups-admin--btn-close");
const admin_share_btn = document.querySelectorAll(".main__groups-admin--btn-share");

// add group section
const add_group_window = document.querySelector(".add_group");
const add_group_close_btn = document.querySelector(".add_group-wrapper--btn");
const add_group_add_btn = document.querySelector(".add_group-btn");
const add_group_create_btn = document.querySelector(".add_group-create--btn");
const add_group_input = document.querySelector(".add_group-input");

// inforamtion windows section
const error_window = document.querySelector(".error__window");
const warning_window = document.querySelector(".warning__window");
const success_window = document.querySelector(".success__window");

// tokens section
const jwt_input = document.getElementById("jwt");
const rjwt_input = document.getElementById("rjwt");

const localStorage = window.localStorage;

localStorage.setItem("jwt", jwt_input.value);
localStorage.setItem("rjwt", rjwt_input.value);

document.body.removeChild(jwt_input);
document.body.removeChild(rjwt_input);

setInterval(()=>{
	newJWT();
}, (1000*60*15 - 30000));
setInterval(()=>{
	newRJWT()
}, (1000*60*60*24*7 - 30000));

// got connection with internet back
window.addEventListener("online", ()=>{
	if (success_window.classList.contains("active")) {
		success_window.classList.remove("active");
	}
	if (error_window.classList.contains("active")) {
		error_window.classList.remove("active");
	}
	let p = warning_window.children[0];
	p.innerHTML = "Восстановленно соединение с интернетом!";
	warning_window.classList.add("active");
	setTimeout(()=>{warning_window.classList.remove("active")}, 4000);
});
// show warning if we lost connection to the internet
window.addEventListener("offline", ()=>{
	if (success_window.classList.contains("active")) {
		success_window.classList.remove("active");
	}
	if (error_window.classList.contains("active")) {
		error_window.classList.remove("active");
	}
	let p = warning_window.children[0];
	p.innerHTML = "Потерянно соединение с интернетом!";
	warning_window.classList.add("active");
});

// delete group logic
delete_btn.forEach((button)=>{
	button.addEventListener("click", ()=>{
		// deleting group
		let parent = button.parentNode.parentNode;
		let child = button.parentNode;

		let group_name = child.children[0].innerText;

		// request to api server
		delete_group(parent, child, group_name);
	});
});
// deleting group where you are admin
admin_delete_btn.forEach((button)=>{
	button.addEventListener("click", ()=>{
		if (confirm("Вы уверены, что хотите удалить группу? Так как вы её владелец, она удалится для всех.")) {	
			// deleting group
			let parent = button.parentNode.parentNode;
			let child = button.parentNode;

			let group_name = child.children[0].innerText;

			// request to api server
			delete_group(parent, child, group_name);
		}
	});
});
// copy link to clipboard
admin_share_btn.forEach((button)=>{
	button.addEventListener("click", ()=>{
		let link = button.parentNode.parentNode.children[0].dataset.link;

		navigator.clipboard.writeText(link).then(()=>{
			showSuccess("Ссылка-приглашение на вступление в группу скопирована в буфер обмена!");
		}, ()=>{
			showError("Ссылка-приглашение на вступление в группу не была скопирована в буфер обмена!");
		});
	});
});

// open group creation window
add_group_create_btn.addEventListener("click", ()=>{
	add_group_input.focus();
	add_group_window.classList.add("active");
});
// close group creation window
add_group_close_btn.addEventListener("click", ()=>{
	add_group_input.value = "";
	add_group_window.classList.remove("active");
});
// create group logic
add_group_add_btn.addEventListener("click", ()=>{
	if (add_group_input.value != "") {	
		let group_name = add_group_input.value.replaceAll(/[\r\n\t]/g, "");
		add_group_input.value = "";
		add_group_input.focus();

		group_exists = false;

		let text_items = document.querySelectorAll(".main__group-item--text");
		text_items.forEach((item)=>{
			if (item.innerText == group_name) {
				showError("Группа с таким именем уже существует!");
				group_exists = true;
			}
		});

		if (!group_exists) {
			parents = [];
			childs = [];
			parent1 = document.querySelector(".main__groups-list");
			parent2 = document.querySelector(".main__groups-admin");

			item_list = document.createElement("div");
			item_link = document.createElement("a");
			item_list_button = document.createElement("button");
			item_list_img = document.createElement("img");

			item_admin = document.createElement("div");
			item_admin_p = document.createElement("p");
			item_admin_button = document.createElement("button");
			item_admin_img_share = document.createElement("img");
			item_admin_img_close = document.createElement("img");

			item_list.classList.add("main__group-item");
			item_link.classList.add("main__group-item--text");
			item_list_button.classList.add("main__group-item--btn");
			item_list_img.classList.add("main__group-item--img");

			item_list_img.src = "/static/dashboard/img/close.svg";
			item_list_img.style.width = "20px";
			item_list_img.style.height = "20px";

			item_link.innerText = group_name;

			item_list_button.appendChild(item_list_img);
			item_list.appendChild(item_link);
			item_list.appendChild(item_list_button);

			item_admin.classList.add("main__group-item");
			item_admin_p.classList.add("main__group-item--text");
			item_admin_p.classList.add("main__groups-admin--text");
			item_admin_button.classList.add("main__groups-admin--btns");
			item_admin_img_share.classList.add("main__group-item--img");
			item_admin_img_close.classList.add("main__group-item--img");
			item_admin_img_share.classList.add("main__groups-admin--btn-share");
			item_admin_img_close.classList.add("main__groups-admin--btn-close");

			item_admin_img_share.src = "/static/dashboard/img/share.svg";
			item_admin_img_share.style.width = "20px";
			item_admin_img_share.style.height = "20px";

			item_admin_img_close.src = "/static/dashboard/img/close.svg";
			item_admin_img_close.style.width = "20px";
			item_admin_img_close.style.height = "20px";

			item_admin_p.innerText = group_name;

			item_admin_button.appendChild(item_admin_img_share);
			item_admin_button.appendChild(item_admin_img_close);
			item_admin.appendChild(item_admin_p);
			item_admin.appendChild(item_admin_button);

			item_admin_img_share.addEventListener("click", ()=>{
				let link = item_admin_img_share.parentNode.parentNode.children[0].dataset.link;

				navigator.clipboard.writeText(link).then(()=>{
					showSuccess("Ссылка-приглашение на вступление в группу скопирована в буфер обмена!");
				}, ()=>{
					showError("Ссылка-приглашение на вступление в группу не была скопирована в буфер обмена!");
				});
			});

			childs[0] = item_list;
			childs[1] = item_admin;

			parents[0] = parent1;
			parents[1] = parent2;

			add_group(parents, childs, group_name);
		}
	} else {
		showError("Имя группы не может быть пустым!");
	}
});

function add_group(parents, childs, group_name) {
	fetch("https://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "add_group",
			"data": {
				"jwt": localStorage.getItem("jwt"),
				"group_name": encodeURI(group_name)
			}
		})
	})
	.then((response) => {
		let status = response.status;

		if (status == 200) {
			return response.json();
		} else if (status == 400) { // bad request
			document.location.reload();
		} else if (status == 401) { // unathorized
			newJWT();
			add_group(parents, childs, group_name);
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробую выполнить данную операцию позже.");
			setTimeout(()=>{
				// try again with 5s delay
				add_group(parents, childs, group_name);
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {
			console.log(childs[0].children[0])
			childs[0].children[0].href = "https://catalog.cc/select/"+data.data.select_link;
			childs[1].children[0].setAttribute("data-link", data.data.invite_link);
			parents[0].appendChild(childs[0]);
			parents[1].appendChild(childs[1]);

			showSuccess("Успешно создана группа!");
		} else {
			showError(data.data.error);
		}
	})
	.catch((error) => {
		showError("Что-то пошло не так, попробуйте перезагрузить страницу!");
	});
}

function delete_group(parent, child, group_name) {
	fetch("https://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "delete_group",
			"data": {
				"jwt": localStorage.getItem("jwt"),
				"group_name": encodeURI(group_name)
			}
		})
	})
	.then((response) => {
		let status = response.status;

		if (status == 200) {
			return response.json();
		} else if (status == 400) { // bad request
			document.location.reload();
		} else if (status == 401) { // unathorized
			newJWT();
			delete_group(parent, child, group_name);
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробую выполнить данную операцию позже.");
			setTimeout(()=>{
				// try again with 5s delay
				delete_group(parent, child, group_name);
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {	
			parent.removeChild(child);

			showSuccess("Успешно удаленно элемент!");
		} else {
			showError(data.data.error);
		}
	})
	.catch((error) => {
		showError("Что-то пошло не так, попробуйте перезагрузить страницу!");
	});
}

function newRJWT() {
	fetch("https://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "new_rjwt",
			"data": {
				"rjwt": localStorage.getItem("rjwt")
			}
		})
	})
	.then((response) => {
		let status = response.status;

		if (status == 200) {
			return response.json();
		} else if (status == 400) { // bad request
			document.location.reload();
		} else if (status == 401) { // unathorized
			document.location.reload();
		} else if (status == 500) { // internal server error
			newRJWT();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробуйте воспользоваться им позже.");
			setTimeout(()=>{
				newRJWT();
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {
			localStorage.setItem("jwt", data.data.jwt);
			localStorage.setItem("rjwt", data.data.rjwt);
		}
	})
	.catch((error) => {});
}

function newJWT() {
	fetch("https://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "new_jwt",
			"data": {
				"jwt": localStorage.getItem("jwt")
			}
		})
	})
	.then((response) => {
		let status = response.status;

		if (status == 200) {
			return response.json();
		} else if (status == 400) { // bad request
			document.location.reload();
		} else if (status == 401) { // unathorized
			newRJWT();
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробуйте воспользоваться им позже.");
			setTimeout(()=>{
				newJWT();
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {
			localStorage.setItem("jwt", data.data.jwt);
		}
	})
	.catch((error) => {});
}

function showError(text) {
	if (success_window.classList.contains("active")) {
		success_window.classList.remove("active");
	}
	if (warning_window.classList.contains("active")) {
		warning_window.classList.remove("active");
	}
	let p = error_window.children[0];
	p.innerHTML = text;
	error_window.classList.add("active");
	setTimeout(()=>{error_window.classList.remove("active")}, 10000);
}

function showWarning(text) {
	if (success_window.classList.contains("active")) {
		success_window.classList.remove("active");
	}
	if (error_window.classList.contains("active")) {
		error_window.classList.remove("active");
	}
	let p = warning_window.children[0];
	p.innerHTML = text;
	warning_window.classList.add("active");
	setTimeout(()=>{warning_window.classList.remove("active")}, 10000);
}

function showSuccess(text) {
	if (error_window.classList.contains("active")) {
		error_window.classList.remove("active");
	}
	if (warning_window.classList.contains("active")) {
		warning_window.classList.remove("active");
	}
	let p = success_window.children[0];
	p.innerHTML = text;
	success_window.classList.add("active");
	setTimeout(()=>{success_window.classList.remove("active")}, 10000);
}

// xss prevention
String.prototype.escape = function() {
	var tagsToReplace = {
		'&': '&amp;',
		'<': '&lt;',
		'>': '&gt;',
		'"': '&quot;',
		'\'': '&apos;',
	};
	return this.replace(/[&<>]/g, function(tag) {
		return tagsToReplace[tag] || tag;
	});
};