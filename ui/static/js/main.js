var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

var themeToggle = document.querySelector("[data-theme-toggle]");

function persistTheme(theme) {
	var maxAge = 60 * 60 * 24 * 365;
	document.cookie = "theme=" + theme + "; Max-Age=" + maxAge + "; Path=/; SameSite=Lax";
}

function applyTheme(theme) {
	document.documentElement.setAttribute("data-theme", theme);
	persistTheme(theme);
	if (themeToggle) {
		themeToggle.textContent = theme === "dark" ? "Light mode" : "Dark mode";
		themeToggle.setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
	}
}

if (themeToggle) {
	var currentTheme = document.documentElement.getAttribute("data-theme") || "light";
	themeToggle.addEventListener("click", function () {
		var nextTheme = document.documentElement.getAttribute("data-theme") === "dark" ? "light" : "dark";
		applyTheme(nextTheme);
	});
}
