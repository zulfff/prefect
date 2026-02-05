// ========================================
// Prefect File Explorer - Modular Components
// ========================================

// ===== 1. State =====
let currentPath = ""; // empty = Home

// ===== 2. DOM Elements =====
const shortcutsEl = document.getElementById("shortcuts");
const drivesEl = document.getElementById("drives");
const filesEl = document.getElementById("files");
const breadcrumbsEl = document.getElementById("breadcrumbs");
const closeBtnEl = document.getElementById("close-btn");

// ===== 3. Icon Mapping =====
const SIDEBAR_ICONS = {
  Home: "icon/home.svg",
  Documents: "icon/docs.svg",
  Downloads: "icon/download.svg",
  Media: "icon/media.svg",
};

const DRIVE_ICON = "icon/mnt.svg";

// ===== 4. Utility Functions =====

/**
 * Get the display name from a mount point path
 * e.g., "/mnt/HDD" -> "HDD", "/" -> "Root"
 */
function getDriveName(mountPath) {
  if (mountPath === "/") return "Root";
  const parts = mountPath.split("/").filter(Boolean);
  return parts[parts.length - 1] || mountPath;
}

/**
 * Get icon source for a sidebar item
 */
function getSidebarIcon(name) {
  return SIDEBAR_ICONS[name] || "icon/folder.svg";
}

/**
 * Get icon source for a file/folder
 */
function getFileIcon(file) {
  const iconName = file.icon || (file.is_dir ? 'folder' : 'file');

  if (file.is_dir) {
    if (iconName === 'folder') {
      return "icon/folder.svg";
    }
    return `folder_icon/${iconName}.svg`;
  }

  if (iconName === 'file') {
    return "icon/file.svg";
  }

  return `icons/${iconName}.svg`;
}

/**
 * Format path for display (replace home with ~)
 */
function formatPathForDisplay(path) {
  if (!path) return "~/";
  // Try to detect home directory pattern and replace with ~
  const homeMatch = path.match(/^(\/home\/[^/]+)/);
  if (homeMatch) {
    return path.replace(homeMatch[1], "~") + "/";
  }
  return path + "/";
}

// ===== 5. Component Functions =====

/**
 * Create a sidebar item element
 */
function SidebarItem(name, path, iconSrc) {
  const div = document.createElement("div");
  div.className = "sidebar-item";
  div.innerHTML = `
    <img src="${iconSrc}" alt="${name}" />
    <span>${name}</span>
  `;
  div.onclick = () => {
    currentPath = path;
    loadExplorer();
  };
  return div;
}

/**
 * Create a file/folder card element
 */
function FileCard(file) {
  const div = document.createElement("div");
  div.className = "file-card";
  const iconSrc = getFileIcon(file);

  div.innerHTML = `
    <img src="${iconSrc}" alt="${file.is_dir ? 'Folder' : 'File'}" />
    <span>${file.name}</span>
  `;

  if (file.is_dir) {
    div.onclick = () => {
      currentPath = file.path;
      loadExplorer();
    };
  }

  return div;
}

/**
 * Create breadcrumb navigation elements
 */
function PathBreadcrumbs(path) {
  const container = document.createDocumentFragment();

  if (!path) {
    // Show just "~/" for home
    const segment = document.createElement("span");
    segment.className = "breadcrumb-segment";
    segment.textContent = "~/";
    container.appendChild(segment);
    return container;
  }

  // Parse path into segments
  const displayPath = formatPathForDisplay(path);
  const isHomePath = displayPath.startsWith("~");

  if (isHomePath) {
    // Handle home-relative paths
    const pathWithoutHome = displayPath.slice(1); // Remove ~
    const segments = pathWithoutHome.split("/").filter(Boolean);

    // Add ~ segment (clickable to go home)
    const homeSegment = document.createElement("span");
    homeSegment.className = "breadcrumb-segment";
    homeSegment.textContent = "~";
    homeSegment.onclick = () => {
      currentPath = "";
      loadExplorer();
    };
    container.appendChild(homeSegment);

    // Build paths for each segment
    const homeMatch = path.match(/^(\/home\/[^/]+)/);
    let buildPath = homeMatch ? homeMatch[1] : "";

    segments.forEach((segment, index) => {
      // Add separator
      const sep = document.createElement("span");
      sep.className = "breadcrumb-separator";
      sep.textContent = "/";
      container.appendChild(sep);

      buildPath += "/" + segment;
      const segmentPath = buildPath;

      const segmentEl = document.createElement("span");
      segmentEl.className = "breadcrumb-segment";
      segmentEl.textContent = segment;
      segmentEl.onclick = () => {
        currentPath = segmentPath;
        loadExplorer();
      };
      container.appendChild(segmentEl);
    });

    // Add trailing slash
    const trailingSep = document.createElement("span");
    trailingSep.className = "breadcrumb-separator";
    trailingSep.textContent = "/";
    container.appendChild(trailingSep);

  } else {
    // Handle absolute paths (drives, etc.)
    const segments = path.split("/").filter(Boolean);
    let buildPath = "";

    // Add root segment
    const rootSegment = document.createElement("span");
    rootSegment.className = "breadcrumb-segment";
    rootSegment.textContent = "/";
    rootSegment.onclick = () => {
      currentPath = "/";
      loadExplorer();
    };
    container.appendChild(rootSegment);

    segments.forEach((segment) => {
      buildPath += "/" + segment;
      const segmentPath = buildPath;

      const segmentEl = document.createElement("span");
      segmentEl.className = "breadcrumb-segment";
      segmentEl.textContent = segment;
      segmentEl.onclick = () => {
        currentPath = segmentPath;
        loadExplorer();
      };
      container.appendChild(segmentEl);

      const sep = document.createElement("span");
      sep.className = "breadcrumb-separator";
      sep.textContent = "/";
      container.appendChild(sep);
    });
  }

  return container;
}

// ===== 6. Render Functions =====

/**
 * Render sidebar shortcuts (Home, Documents, Downloads, Media)
 */
function renderShortcuts(sidebar) {
  shortcutsEl.innerHTML = "";

  sidebar.forEach((item) => {
    const iconSrc = getSidebarIcon(item.directory_name);
    const sidebarItem = SidebarItem(item.directory_name, item.directory_path, iconSrc);
    shortcutsEl.appendChild(sidebarItem);
  });
}

/**
 * Render mounted drives
 */
function renderDrives(drives) {
  drivesEl.innerHTML = "";

  drives.forEach((drive) => {
    const driveName = getDriveName(drive.mount_point);
    const driveItem = SidebarItem(driveName, drive.mount_point, DRIVE_ICON);
    drivesEl.appendChild(driveItem);
  });
}

/**
 * Render files and folders grid
 */
function renderFiles(files) {
  filesEl.innerHTML = "";

  if (!files || files.length === 0) {
    const emptyState = document.createElement("div");
    emptyState.className = "empty-state";
    emptyState.innerHTML = `
      <span>This folder is empty</span>
    `;
    filesEl.appendChild(emptyState);
    return;
  }

  // Sort: folders first, then files, both alphabetically
  const sorted = [...files].sort((a, b) => {
    if (a.is_dir && !b.is_dir) return -1;
    if (!a.is_dir && b.is_dir) return 1;
    return a.name.localeCompare(b.name);
  });

  sorted.forEach((file) => {
    const fileCard = FileCard(file);
    filesEl.appendChild(fileCard);
  });
}

/**
 * Render breadcrumb path bar
 */
function renderBreadcrumbs(path) {
  breadcrumbsEl.innerHTML = "";
  const breadcrumbs = PathBreadcrumbs(path);
  breadcrumbsEl.appendChild(breadcrumbs);
}

// ===== 7. Main Loader =====

function loadExplorer() {
  const url = currentPath
    ? `/api/explorer?path=${encodeURIComponent(currentPath)}`
    : `/api/explorer`;

  fetch(url)
    .then((res) => res.json())
    .then((data) => {
      renderShortcuts(data.sidebar);
      renderDrives(data.drives);
      renderFiles(data.files);
      renderBreadcrumbs(currentPath);
    })
    .catch((err) => {
      console.error("Failed to load explorer:", err);
      filesEl.innerHTML = `
        <div class="empty-state">
          <span>Failed to load files</span>
        </div>
      `;
    });
}

// ===== 8. Event Listeners =====

// Close button - navigate back to index
closeBtnEl.onclick = () => {
  window.location.href = "index.html";
};

// ===== 9. Initial Load =====
loadExplorer();