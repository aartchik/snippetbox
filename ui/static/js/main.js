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

var avatarForm = document.querySelector("[data-avatar-form]");
if (avatarForm) {
	var avatarInput = avatarForm.querySelector('input[type="file"]');
	var avatarButton = avatarForm.querySelector('button[type="submit"]');
	var avatarHint = avatarForm.querySelector("[data-avatar-hint]");

	function updateAvatarFormState() {
		var hasFile = avatarInput && avatarInput.files && avatarInput.files.length > 0;
		if (avatarButton) {
			avatarButton.disabled = !hasFile;
		}
		if (avatarHint) {
			avatarHint.classList.toggle("is-error", !hasFile);
			avatarHint.textContent = hasFile ? "Ready to upload." : "Choose a JPG or PNG before saving.";
		}
	}

	if (avatarInput) {
		avatarInput.addEventListener("change", updateAvatarFormState);
	}

	avatarForm.addEventListener("submit", function (event) {
		if (!avatarInput || !avatarInput.files || avatarInput.files.length === 0) {
			event.preventDefault();
			updateAvatarFormState();
		}
	});

	updateAvatarFormState();
}
