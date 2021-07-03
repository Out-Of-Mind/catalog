// link section
const first_link = document.querySelector(".first-link");
const second_link = document.querySelector(".second-link");

// inputs section
const inputs = document.querySelectorAll(".form__register-input");
const login_input = document.getElementById("login");
const email_input = document.getElementById("email");
const password_input = document.getElementById("password");
const repassword_input = document.getElementById("repassword");
const checkbox_label = document.querySelector(".form__register-checkbox");
const submit = document.querySelector(".form__register-submit");

// error window
const error_window = document.querySelector(".error__window");

let errors = {};

submit.addEventListener("click", (event)=>{
	validate("login", (login_input.value != ""), login_input, "");
	validate("email", validateEmail(email_input.value), email_input, "Введён не верный эмейл!");
	validate("password", validatePassword(password_input.value), password_input, "Пароль должен содержать: одну маленькую букву, одну цифру, как минимум 8 символов!");
	validate("repassword", (repassword_input.value == password_input.value), repassword_input, "Пароли не совпадают!");

	if (errors.login || errors.email || errors.password || errors.repassword) {
		event.preventDefault();
		showError("Все поля должны быть корректно заполнены!");
	}
});

first_link.addEventListener("click", ()=>{
	if (second_link.classList.contains("active")) {
		second_link.classList.remove("active");
		first_link.classList.add("active");
	}
});

second_link.addEventListener("click", ()=>{
	if (first_link.classList.contains("active")) {
		first_link.classList.remove("active");
		second_link.classList.add("active");
	}
});

login_input.addEventListener("input", ()=>{
	validate("login", (login_input.value != ""), login_input, "");
});

email_input.addEventListener("input", ()=>{
	validate("email", validateEmail(email_input.value), email_input, "Введён не верный эмейл!");
});

password_input.addEventListener("input", ()=>{
	validate("password", validatePassword(password_input.value), password_input, "Пароль должен содержать: одну маленькую букву, одну цифру, как минимум 8 символов!");
	validate("repassword", (repassword_input.value == password_input.value), repassword_input, "Пароли не совпадают!");
});

repassword_input.addEventListener("input", ()=>{
	validate("repassword", (repassword_input.value == password_input.value), repassword_input, "Пароли не совпадают!");
});

checkbox_label.addEventListener("click", ()=>{
	let checkbox = checkbox_label.children[0];
	if (checkbox.checked) {
		password_input.setAttribute('type','text');
		repassword_input.setAttribute('type','text');
	} else {
		password_input.setAttribute('type','password');
		repassword_input.setAttribute('type','password');
	}
});

function validate(where, statement, input, error) {
	if (input.value != "") {
		let label = document.getElementById(where+"-label");
		let img = document.createElement("img");

		label.innerHTML = "";
		if (statement == true) {
			errors[where] = false;
			img.src = "/static/register/img/valid.svg";
			hideError();
		} else {
			errors[where] = true;
			img.src = "/static/register/img/invalid.svg";
			showError(error);
		}
		label.appendChild(img);
	} else {
		let label = document.getElementById(where+"-label");
		label.innerHTML = ""
	}
}

function validateEmail(email) {
	return /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(email);
}

function validatePassword(password) {
	return /^(?=.*[a-z])(?=.*[0-9])(?=.{8,})/.test(password);
}

function showError(text) {
	let p = error_window.children[0];
	p.innerHTML = text;
	error_window.classList.add("active");
}

function hideError() {
	error_window.classList.remove("active");
}