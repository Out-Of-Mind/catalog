// header section
const search_input = document.querySelector(".header__search-input");
const search_btn = document.querySelector(".header__search-btn");
const add_window_btn = document.querySelector(".header__add-btn");

// main section
const delete_btn = document.querySelectorAll(".main__group-item--btn");
const main = document.querySelector(".main");

// add item section
const add_window = document.querySelector(".add");
const add_window_close_btn = document.querySelector(".add__btn-close");
const add_window_add_btn = document.querySelector(".add__btn-add");
const add_window_input_category = document.querySelector(".add__input-category");
const add_window_input_description = document.querySelector(".add__input-description");

// menu section
const menu_btn = document.querySelector(".menu-btn");
const menu_img = document.querySelector(".menu-img");
const menu_window = document.querySelector(".menu__window");
const menu_window_add = document.querySelector(".menu__window-add");
const menu_window_delete = document.querySelector(".menu__window-delete");
const menu_window_add_close = document.querySelector(".menu__window-add__btn-close");
const menu_window_delete_close = document.querySelector(".menu__window-delete__btn-close")
const menu_window_add_btn = document.querySelector(".menu__window-btn--add");
const menu_window_delete_btn = document.querySelector(".menu__window-btn--delete");
const menu_window_add_category = document.querySelector(".menu__window-add__btn");
const menu_window_delete_category = document.querySelector(".menu__window-delete__btn");
const menu_window_add_input_category = document.querySelector(".menu__window-add__input-category");
const menu_window_delete_input_category = document.querySelector(".menu__window-delete__input-category");

// search section
const search_div = document.querySelector(".search");
const search_output = document.querySelector(".search__output");
const search_close = document.querySelector(".search__nav-btn");

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

// search logic
search_btn.addEventListener("click", ()=>{
	if (search_input.value == "") {
		showError("Не введено никаких данных!");
	} else {
		search(search_input.value);
		search_input.value = "";
	}
});
// return to main from search output
search_close.addEventListener("click", ()=>{
	search_output.classList.remove("active");
	search_div.classList.remove("active");
	main.classList.remove("disable");

	search_output.innerHTML = "";
});
// add new item opener
add_window_btn.addEventListener("click", ()=>{
	add_window.classList.add("active");
	add_window_input_category.focus();
});
// add new item logic
add_window_add_btn.addEventListener("click", ()=>{
	// check if all values was entered
	if (add_window_input_category.value == "" || add_window_input_description.value == "") {
		showError("Все поля должны быть заполнены!");
	} else {
		// finding category
		let category_name = add_window_input_category.value;
		let description = (add_window_input_description.value).replaceAll(/[\r\n\t]/g, "").escape();
		let group = document.getElementById((category_name).escape().replaceAll(" ", "_"));

		// chech if group with "category_name" was found
		if (!group) {
			showError("Категории с таким именем не найдено!");
		} else {
			if (error_window.classList.contains("active")) {
				error_window.classList.remove("active");
			}
			// adding new item
			add_window_input_description.value = "";

			let group_item = document.createElement("div");
			let group_item_text = document.createElement("p");
			let group_item_btn = document.createElement("button");
			let group_item_img = document.createElement("img");

			group_item.classList.add("main__group-item");
			group_item_text.classList.add("main__group-item--text");
			group_item_btn.classList.add("main__group-item--btn");
			group_item_img.classList.add("main__group-item--img");

			group_item_img.src = "/static/index/img/close.svg";
			group_item_img.alt = "close";
			group_item_img.style.width = "20px";
			group_item_img.style.height = "20px";
			group_item_text.innerHTML = description;

			group_item_btn.appendChild(group_item_img);
			group_item.appendChild(group_item_text);
			group_item.appendChild(group_item_btn);

			group_item_btn.addEventListener("click", ()=>{
				let child = group_item_btn.parentNode;
				group_item_btn.parentNode.parentNode.removeChild(child);
			});

			// request to api server
			add_item(group, group_item, group_item_text.innerText, category_name);
		}
	}
});
// delete item logic
delete_btn.forEach((button)=>{
	button.addEventListener("click", ()=>{
		// deleting item
		let parent = button.parentNode.parentNode;
		let child = button.parentNode;

		let group_item_text = child.innerText;
		let category_name = parent.children[0].innerText;

		// request to api server
		delete_item(parent, child, group_item_text, category_name);
	});
});
// close add category window
add_window_close_btn.addEventListener("click", ()=>{
	add_window.classList.remove("active");
	add_window_input_category.value = "";
	add_window_input_description.value = "";
});
// menu opener
menu_btn.addEventListener("click", ()=>{
	menu_img.classList.toggle("animate");
	menu_window.classList.toggle("active");
});
// add new category opener
menu_window_add_btn.addEventListener("click", ()=>{
	if (menu_window_delete.classList.contains("active")) {
		menu_window_delete.classList.remove("active");
	}
	menu_window_add.classList.add("active");
	menu_window_add_input_category.focus();
});
// delete category opener
menu_window_delete_btn.addEventListener("click", ()=>{
	if (menu_window_add.classList.contains("active")) {
		menu_window_add.classList.remove("active");
	}
	menu_window_delete.classList.add("active");
	menu_window_delete_input_category.focus();
});
// close add category window
menu_window_add_close.addEventListener("click", ()=>{
	menu_window_add.classList.remove("active");
	menu_window_add_input_category.value = "";
});
// close delete category window
menu_window_delete_close.addEventListener("click", ()=>{
	menu_window_delete.classList.remove("active");
	menu_window_delete_input_category.value = "";
});
// add new category logic
menu_window_add_category.addEventListener("click", ()=>{
	let no_error = true;

	// finding category
	let category_name = (menu_window_add_input_category.value).escape();
	let category_id = (category_name.replaceAll(" ", "_"));
	let category = document.getElementById(category_id);

	if (category_name == "") {
		no_error = false;
		showError("У категории не может быть пустым имя!");
	}
	if (category) {
		// category with this name is already exist
		no_error = false;
		showError("Категория с таким именем уже существует!");
	}
	
	if (no_error) {
		// adding new category
		let category = document.createElement("div");
		let category_title = document.createElement("p");

		let category_id = category_name.replaceAll(" ", "_");

		category.classList.add("main__group");
		category_title.classList.add("main__group-title");

		category_title.innerHTML = category_name;
		category.setAttribute("id", category_id);

		category.appendChild(category_title);

		menu_window_add_input_category.value = "";

		// request to api server
		add_category(main, category, category_name);
	}
});
// delete new category logic
menu_window_delete_category.addEventListener("click", ()=>{
	let no_error = true;

	let category_name = menu_window_delete_input_category.value;

	if (category_name == "") {
		showError("Чтобы удалить категорию должно быть заполненно имя!");
	}

	// finding category
	let category_id = (category_name.replaceAll(" ", "_")).escape();
	console.log(category_id)
	let category = document.getElementById(category_id);
	
	if (category == null) {
		// if no category was found show error
		no_error = false;
		showError("Чтобы удалить категорию нужно ввести её имя!");
	}

	if (no_error) {
		// deleting category
		menu_window_delete_input_category.value = "";

		let parent = category.parentNode;

		// request to api server
		delete_category(parent, category, category_name);
	}
});

function add_category(parent, category, category_name) {
	fetch("http://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "add_category",
			"data": {
				"jwt": localStorage.getItem("jwt"),
				"category_name": encodeURI(category_name)
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
			add_category(parent, category, category_name);
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробую выполнить данную операцию позже.");
			setTimeout(()=>{
				add_category(parent, category, category_name);
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {	
			parent.appendChild(category);

			showSuccess("Успешно дабовленно новую категорию!");
		} else {
			showError(data.data.error);
		}
	})
	.catch((error) => {
		showError("Что-то пошло не так, попробуйте перезагрузить страницу!");
	});
}
function delete_category(parent, category, category_name) {
	fetch("http://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "delete_category",
			"data": {
				"jwt": localStorage.getItem("jwt"),
				"category_name": encodeURI(category_name)
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
			delete_category(parent, category_name);
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробую выполнить данную операцию позже.");
			setTimeout(()=>{
				delete_category(parent, category_name);
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {	
			parent.removeChild(category);

			showSuccess("Успешно удаленно категорию!");
		} else {
			showError(data.data.error);
		}
	})
	.catch((error) => {
		showError("Что-то пошло не так, попробуйте перезагрузить страницу!");
	});
}

function add_item(parent, group_item, group_item_text, category_name) {
	fetch("http://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "add_item",
			"data": {
				"jwt": localStorage.getItem("jwt"),
				"item_name": encodeURI(group_item_text),
				"category_name": encodeURI(category_name)
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
			add_item(parent, group_item, group_item_text, category_name);
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробую выполнить данную операцию позже.");
			setTimeout(()=>{
				// try again with 5s delay
				add_item(parent, group_item, group_item_text, category_name);
			}, 5000);
		}
	})
	.then((data) => {
		if (data.success) {	
			parent.appendChild(group_item);

			showSuccess("Успешно дабовленно новый элемент!");
		} else {
			showError(data.data.error);
		}
	})
	.catch((error) => {
		showError("Что-то пошло не так попробуйте перезагрузить страницу!");
	});
}
function delete_item(parent, child, group_item_text, category_name) {
	fetch("http://api.catalog.cc/", {
		method: "POST",
		body: JSON.stringify({
			"action": "delete_item",
			"data": {
				"jwt": localStorage.getItem("jwt"),
				"item_name": encodeURI(group_item_text),
				"category_name": encodeURI(category_name)
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
			delete_item(parent, child, group_item_text, category_name);
		} else if (status == 500) { // internal server error
			document.location.reload();
		} else if (status == 503) { // service unaviable
			showWarning("Сервис временно недоступен, попробую выполнить данную операцию позже.");
			setTimeout(()=>{
				// try again with 5s delay
				delete_item(parent, child, group_item_text, category_name);
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

function search(str) {
	const main_group = document.querySelectorAll(".main__group");
	// array of strings where we will find key word
	let array = [];
	// array of found items
	let output = [];
	
	main_group.forEach(group=>{
		let inner_strings = [];

		// check if group has any children
		if (group.hasChildNodes()) {
			let children = group.childNodes;
			let length = children.length;

			children.forEach(child=>{
				// nodeType 3 don't needed
				if (child.nodeType != 3) {
					if (child.childElementCount != 0) {
						// if it's item:
						let text = child.children[0].innerText;
						inner_strings.push(text);
					} else {
						// if it's title:
						let text = child.innerText;
						inner_strings.push(text);
					}
				}
			})
		}

		array.push({"innerText": inner_strings});
	});

	// search itself
	array.forEach(element=>{
		let found = false;

		// find logic
		element["innerText"].forEach(string=>{
			if (string.includes(str)) {
				found = true;
			}
		});

		if (found) {
			// if we find key word push title of it's parent into the output array
			output.push(element["innerText"][0]);
		}
	});
	
	if (output.length > 0) {
		// if found someting
		output.forEach(element=>{
			// enumarte output array
			let category_name = element;
			let group = document.getElementById((element).escape().replaceAll(" ", "_")).cloneNode(true);
			
			// copy elements from main to search output
			search_output.appendChild(group);
		});

		// activating search output visibility
		search_output.classList.add("active");
		search_div.classList.add("active");
		main.classList.add("disable");

		showWarning("Здесь вы можете только просматривать карточки!");
	} else {
		// nothing was found
		showError("По вашему запросу ничего не найдено.");
	}
}

function newRJWT() {
	fetch("http://api.catalog.cc/", {
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
	fetch("http://api.catalog.cc/", {
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