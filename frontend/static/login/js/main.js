const first_link = document.querySelector(".first-link");
const second_link = document.querySelector(".second-link");

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