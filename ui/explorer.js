// 1. Grab DOM elements
const filesEl = document.getElementById("files");
const sidebarEl = document.getElementById("sidebar");
const mountEl = document.getElementById("mount");

// 2. Fetch explorer data
fetch("/api/explorer")
  .then(res => res.json())
  .then(data => {
    renderSidebar(data.sidebar);
    renderDrives(data.drives);
    renderFiles(data.files);
  })
  .catch(err => {
    console.error("Failed to load explorer:", err);
  });

// 3. Render sidebar
function renderSidebar(items) {
  sidebarEl.innerHTML = "<h3>Sidebar</h3>";

  items.forEach(item => {
    const div = document.createElement("div");
    div.textContent = item.directory_name;
    sidebarEl.appendChild(div);
  });
}

// 4. Render drives
function renderDrives(drives) {
  mountEl.innerHTML = "<h3>Drives</h3>";

  drives.forEach(drive => {
    const div = document.createElement("div");
    div.textContent = `${drive.device} → ${drive.mount_point}`;
    mountEl.appendChild(div);
  });
}

// 5. Render files
function renderFiles(files) {
  filesEl.innerHTML = "<h3>Files</h3>";

  files.forEach(file => {
    const div = document.createElement("div");

    if (file.is_dir) {
      div.textContent = "📁 " + file.name;
    } else {
      div.textContent = "📄 " + file.name;
    }

    filesEl.appendChild(div);
  });
}
