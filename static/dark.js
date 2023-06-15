const toggleBtn = document.getElementById("toggle-btn");
let darkMode = localStorage.getItem("dark-mode");

const enableDarkMode = () => {
  document.documentElement.classList.add("dark");
  toggleBtn.classList.remove("dark-mode-toggle");
  localStorage.setItem("dark-mode", "enabled");
};

const disableDarkMode = () => {
  document.documentElement.classList.remove("dark");
  toggleBtn.classList.add("dark-mode-toggle");
  localStorage.setItem("dark-mode", "disabled");
};

if (darkMode === "enabled") {
  enableDarkMode(); 
}

if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
  if (!darkMode) {
    enableDarkMode();
  }
}

toggleBtn.addEventListener("click", () => {
  darkMode = localStorage.getItem("dark-mode");
  if (darkMode === "disabled") {
    enableDarkMode();
  } else {
    disableDarkMode();
  }
});
