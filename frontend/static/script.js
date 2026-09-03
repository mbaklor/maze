/**
 * @returns { "system" | "light" | "dark" } the current theme
 */
function getCurrentTheme() {
  return document.querySelector(".theme-chooser input:checked").value;
}

/**
 * @param { "system" | "light" | "dark" } theme to switch to
 */
function setTheme(theme) {
  if (theme === "system") {
    localStorage.removeItem("theme");
  } else {
    localStorage.setItem("theme", theme);
  }
  document.querySelector(`#theme-picker-${theme}`).checked = true;
}

function startup() {
  document.documentElement.style.setProperty("color-scheme", null);
  document.documentElement.style.setProperty("--theme", null);
  setTheme(currentTheme);
}
startup();

function switchTheme() {
  switch (getCurrentTheme()) {
    case "system":
      setTheme("light");
      break;
    case "light":
      setTheme("dark");
      break;
    case "dark":
      setTheme("system");
      break;
  }
}
