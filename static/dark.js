const enableDarkMode = () => {
  document.documentElement.classList.add("dark");
  setCookie("dark-mode", "enabled", 365);
  const toggleBtn = document.getElementById("toggle-btn");
  toggleBtn.checked = true;
};

const disableDarkMode = () => {
  document.documentElement.classList.remove("dark");
  setCookie("dark-mode", "disabled", 365);
};

const setCookie = (cname, cvalue, exdays) => {
  const d = new Date();
  d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
  let expires = "expires=" + d.toUTCString();
  document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/;SameSite=Strict";
}

const getCookie = (cname) => {
  let name = cname + "=";
  let ca = document.cookie.split(';');
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

let darkMode = getCookie("dark-mode");
console.log(darkMode)

if (darkMode === "enabled") {
  enableDarkMode();
}

if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
  if (!darkMode) {
    enableDarkMode();
  }
}

document.body.addEventListener("change", (element) => {
  if (element.target.id === "toggle-btn") {
    darkMode = getCookie("dark-mode");
    if (darkMode === "disabled") {
      enableDarkMode();
      element.target.checked = true;
    } else {
      disableDarkMode();
      element.target.checked = false;
    }
  }
});
