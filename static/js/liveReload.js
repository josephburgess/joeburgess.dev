(function () {
  const checkForReload = () => {
    fetch("/check-reload")
      .then((response) => response.text())
      .then((data) => {
        if (data === "reload") {
          console.log("Reloading page...");
          window.location.reload();
        }
      })
      .catch((err) => console.error("Error checking for reload:", err));
  };

  setInterval(checkForReload, 1000);
  console.log("Live reload enabled");
})();
