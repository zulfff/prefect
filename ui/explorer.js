// ========================================
// Prefect File Explorer - Modular Components
// ========================================

// ===== 1. State =====
let currentPath = ""; // empty = Home
let homePath = ""; // Will be updated from API
let selectedFile = null; // Currently selected file data
let clipboard = {
  items: [], // Array of file data objects
  operation: null, // 'copy' or 'cut'
};

// ===== 2. DOM Elements =====
const shortcutsEl = document.getElementById("shortcuts");
const drivesEl = document.getElementById("drives");
const filesEl = document.getElementById("files");
const breadcrumbsEl = document.getElementById("breadcrumbs");
const closeBtnEl = document.getElementById("close-btn");
const contextMenuEl = document.getElementById("context-menu");
const renameModalEl = document.getElementById("rename-modal");
const deleteModalEl = document.getElementById("delete-modal");

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
  if (!path || path === homePath) return "~";

  if (!homePath) {
    // Fallback until homePath is loaded
    const homeMatch = path.match(/^(\/home\/[^/]+)/);
    if (homeMatch) return path.replace(homeMatch[1], "~");
    return path;
  }

  if (path.startsWith(homePath)) {
    return "~" + path.slice(homePath.length);
  }
  return path;
}

/**
 * Convert absolute path to relative path (for API calls)
 * Returns paths relative to home without leading slash, or absolute paths for drives
 */
function toRelativePath(absPath) {
  if (!absPath) return "";
  // Remove home directory prefix
  const homeMatch = absPath.match(/^\/home\/[^/]+(.*)$/);
  if (homeMatch) {
    // Remove leading slash to make it a proper relative path
    const relativePart = homeMatch[1] || "";
    return relativePart.startsWith('/') ? relativePart.substring(1) : relativePart;
  }
  // For non-home paths (like /mnt/HDD), return as-is (absolute)
  return absPath;
}

// ===== 5. API Functions =====

async function apiDelete(path, isDir) {
  const endpoint = isDir ? '/api/explorer/delete' : '/api/explorer/delete';
  const relativePath = toRelativePath(path);
  const url = `${endpoint}?path=${encodeURIComponent(relativePath)}`;

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return true;
}

async function apiRename(path, newName) {
  const relativePath = toRelativePath(path);
  const url = `/api/explorer/rename?path=${encodeURIComponent(relativePath)}&name=${encodeURIComponent(newName)}`;

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return true;
}

async function apiCopy(srcPath, dstDir) {
  const relativeSrc = toRelativePath(srcPath);
  const relativeDst = toRelativePath(dstDir);
  const url = `/api/explorer/copy?src=${encodeURIComponent(relativeSrc)}&dst=${encodeURIComponent(relativeDst)}`;

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return true;
}

async function apiCut(srcPath, dstDir) {
  const relativeSrc = toRelativePath(srcPath);
  const relativeDst = toRelativePath(dstDir);
  const url = `/api/explorer/cut?src=${encodeURIComponent(relativeSrc)}&dst=${encodeURIComponent(relativeDst)}`;

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return true;
}

// ===== 6. Clipboard Operations =====

function copyToClipboard(file) {
  clipboard.items = [file];
  clipboard.operation = 'copy';
  updateCutVisuals();
}

function cutToClipboard(file) {
  clipboard.items = [file];
  clipboard.operation = 'cut';
  updateCutVisuals();
}

function clearClipboard() {
  clipboard.items = [];
  clipboard.operation = null;
  updateCutVisuals();
}

function updateCutVisuals() {
  // Remove all cut classes
  document.querySelectorAll('.file-card.cut').forEach(el => {
    el.classList.remove('cut');
  });

  // Add cut class to items in clipboard if operation is 'cut'
  if (clipboard.operation === 'cut') {
    clipboard.items.forEach(item => {
      const card = document.querySelector(`.file-card[data-path="${item.path}"]`);
      if (card) {
        card.classList.add('cut');
      }
    });
  }
}

async function pasteFromClipboard() {
  if (clipboard.items.length === 0) return;

  const destPath = currentPath || "";

  try {
    for (const item of clipboard.items) {
      if (clipboard.operation === 'copy') {
        await apiCopy(item.path, destPath);
      } else if (clipboard.operation === 'cut') {
        await apiCut(item.path, destPath);
      }
    }

    if (clipboard.operation === 'cut') {
      clearClipboard();
    }

    loadExplorer();
  } catch (err) {
    console.error("Paste failed:", err);
    alert("Failed to paste: " + err.message);
  }
}

function downloadFile(file) {
  if (file.is_dir) {
    alert("Cannot download directories");
    return;
  }

  try {
    // Build the download URL with the file path
    const downloadUrl = `/api/download?path=${encodeURIComponent(file.path)}`;
    
    // Create a temporary link and trigger download
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = file.name;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  } catch (err) {
    console.error("Download failed:", err);
    alert("Failed to download: " + err.message);
  }
}

// ===== 7. Context Menu =====

function showContextMenu(x, y, file) {
  selectedFile = file;

  // Position the menu
  contextMenuEl.style.left = `${x}px`;
  contextMenuEl.style.top = `${y}px`;

  // Enable/disable menu items based on context
  const copyBtn = document.getElementById('ctx-copy');
  const cutBtn = document.getElementById('ctx-cut');
  const pasteBtn = document.getElementById('ctx-paste');
  const renameBtn = document.getElementById('ctx-rename');
  const deleteBtn = document.getElementById('ctx-delete');
  const downloadBtn = document.getElementById('ctx-download');

  if (file) {
    copyBtn.classList.remove('disabled');
    cutBtn.classList.remove('disabled');
    renameBtn.classList.remove('disabled');
    deleteBtn.classList.remove('disabled');
    // Only enable download for files, not directories
    if (!file.is_dir) {
      downloadBtn.classList.remove('disabled');
    } else {
      downloadBtn.classList.add('disabled');
    }
  } else {
    copyBtn.classList.add('disabled');
    cutBtn.classList.add('disabled');
    renameBtn.classList.add('disabled');
    deleteBtn.classList.add('disabled');
    downloadBtn.classList.add('disabled');
  }

  // Enable paste only if clipboard has items
  if (clipboard.items.length > 0) {
    pasteBtn.classList.remove('disabled');
  } else {
    pasteBtn.classList.add('disabled');
  }

  // Show menu
  contextMenuEl.classList.add('visible');

  // Adjust position if menu goes off screen
  const rect = contextMenuEl.getBoundingClientRect();
  if (rect.right > window.innerWidth) {
    contextMenuEl.style.left = `${window.innerWidth - rect.width - 10}px`;
  }
  if (rect.bottom > window.innerHeight) {
    contextMenuEl.style.top = `${window.innerHeight - rect.height - 10}px`;
  }
}

function hideContextMenu() {
  contextMenuEl.classList.remove('visible');
}

// ===== 8. Modal Functions =====

function showRenameModal(file) {
  const input = document.getElementById('rename-input');
  input.value = file.name;
  renameModalEl.classList.add('visible');

  // Select filename without extension for files
  setTimeout(() => {
    input.focus();
    if (!file.is_dir && file.name.includes('.')) {
      const lastDot = file.name.lastIndexOf('.');
      input.setSelectionRange(0, lastDot);
    } else {
      input.select();
    }
  }, 50);
}

function hideRenameModal() {
  renameModalEl.classList.remove('visible');
}

async function confirmRename() {
  const input = document.getElementById('rename-input');
  const newName = input.value.trim();

  if (!newName || !selectedFile) {
    hideRenameModal();
    return;
  }

  if (newName === selectedFile.name) {
    hideRenameModal();
    return;
  }

  try {
    await apiRename(selectedFile.path, newName);
    hideRenameModal();
    loadExplorer();
  } catch (err) {
    console.error("Rename failed:", err);
    alert("Failed to rename: " + err.message);
  }
}

function showDeleteModal(file) {
  const message = document.getElementById('delete-message');
  message.textContent = `Are you sure you want to delete "${file.name}"?`;
  deleteModalEl.classList.add('visible');
}

function hideDeleteModal() {
  deleteModalEl.classList.remove('visible');
}

async function confirmDelete() {
  if (!selectedFile) {
    hideDeleteModal();
    return;
  }

  try {
    await apiDelete(selectedFile.path, selectedFile.is_dir);
    hideDeleteModal();
    selectedFile = null;
    loadExplorer();
  } catch (err) {
    console.error("Delete failed:", err);
    alert("Failed to delete: " + err.message);
  }
}

// ===== 9. Component Functions =====

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
  div.setAttribute('data-path', file.path);
  const iconSrc = getFileIcon(file);

  div.innerHTML = `
    <img src="${iconSrc}" alt="${file.is_dir ? 'Folder' : 'File'}" />
    <span>${file.name}</span>
  `;

  // Left click - navigate into folders or select files
  div.onclick = (e) => {
    e.stopPropagation();

    // Clear previous selection
    document.querySelectorAll('.file-card.selected').forEach(el => {
      el.classList.remove('selected');
    });

    // Select this card
    div.classList.add('selected');
    selectedFile = file;

    // Double-click to navigate
    if (file.is_dir) {
      if (e.detail === 2) { // Double click
        currentPath = file.path;
        loadExplorer();
      }
    }
  };

  // Right click - show context menu
  div.oncontextmenu = (e) => {
    e.preventDefault();
    e.stopPropagation();

    // Select the card
    document.querySelectorAll('.file-card.selected').forEach(el => {
      el.classList.remove('selected');
    });
    div.classList.add('selected');

    showContextMenu(e.clientX, e.clientY, file);
  };

  // Check if this file is in the cut clipboard
  if (clipboard.operation === 'cut' && clipboard.items.some(item => item.path === file.path)) {
    div.classList.add('cut');
  }

  return div;
}

/**
 * Create breadcrumb navigation elements
 */
function PathBreadcrumbs(path) {
  const container = document.createDocumentFragment();
  const displayPath = formatPathForDisplay(path);
  const isHome = displayPath.startsWith("~");

  // Split path into parts and remove empty ones
  const parts = displayPath.split("/").filter(Boolean);

  // If initial path was empty or resolved to just "~"
  if (parts.length === 0 || (parts.length === 1 && parts[0] === "~")) {
    const rootSeg = document.createElement("span");
    rootSeg.className = "breadcrumb-segment";
    rootSeg.textContent = isHome ? "~" : "Root";
    rootSeg.onclick = () => {
      currentPath = isHome ? "" : "/";
      loadExplorer();
    };
    container.appendChild(rootSeg);

    const sep = document.createElement("span");
    sep.className = "breadcrumb-separator";
    sep.textContent = "/";
    container.appendChild(sep);
    return container;
  }

  if (isHome) {
    // Home paths: "~", "Downloads", etc.
    // parts[0] is "~"
    let fullPath = homePath;

    parts.forEach((part, index) => {
      const segEl = document.createElement("span");
      segEl.className = "breadcrumb-segment";
      segEl.textContent = part;

      let targetPath;
      if (part === "~") {
        targetPath = homePath;
      } else {
        fullPath = fullPath.endsWith("/") ? fullPath + part : fullPath + "/" + part;
        targetPath = fullPath;
      }

      segEl.onclick = () => {
        currentPath = targetPath;
        loadExplorer();
      };
      container.appendChild(segEl);

      const sep = document.createElement("span");
      sep.className = "breadcrumb-separator";
      sep.textContent = "/";
      container.appendChild(sep);
    });
  } else {
    // Absolute paths: "/", "mnt", "HDD", etc.
    // Render the root segment first
    const rootSeg = document.createElement("span");
    rootSeg.className = "breadcrumb-segment";
    rootSeg.textContent = "Root";
    rootSeg.onclick = () => {
      currentPath = "/";
      loadExplorer();
    };
    container.appendChild(rootSeg);

    const rootSep = document.createElement("span");
    rootSep.className = "breadcrumb-separator";
    rootSep.textContent = "/";
    container.appendChild(rootSep);

    let fullPath = "";
    parts.forEach((part) => {
      fullPath += "/" + part;
      const targetPath = fullPath;

      const segEl = document.createElement("span");
      segEl.className = "breadcrumb-segment";
      segEl.textContent = part;
      segEl.onclick = () => {
        currentPath = targetPath;
        loadExplorer();
      };
      container.appendChild(segEl);

      const sep = document.createElement("span");
      sep.className = "breadcrumb-separator";
      sep.textContent = "/";
      container.appendChild(sep);
    });
  }

  return container;
}

// ===== 10. Render Functions =====

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

// ===== 11. Main Loader =====

function loadExplorer() {
  const url = currentPath
    ? `/api/explorer?path=${encodeURIComponent(currentPath)}`
    : `/api/explorer`;

  fetch(url)
    .then((res) => res.json())
    .then((data) => {
      if (data.sidebar && data.sidebar.length > 0) {
        homePath = data.sidebar[0].directory_path;
      }
      renderShortcuts(data.sidebar);
      renderDrives(data.drives);
      renderFiles(data.files);
      renderBreadcrumbs(currentPath);
      selectedFile = null;
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

// ===== 12. Event Listeners =====

// Close button - navigate back to index
closeBtnEl.onclick = () => {
  // Check if running in iframe
  if (window.parent && window.parent !== window) {
    window.parent.postMessage('close-app', '*');
  } else {
    window.location.href = "index.html";
  }
};

// Hide context menu on click outside
document.addEventListener('click', (e) => {
  if (!contextMenuEl.contains(e.target)) {
    hideContextMenu();
  }

  // Deselect if clicking on empty space in files grid
  if (e.target === filesEl || e.target.classList.contains('empty-state')) {
    document.querySelectorAll('.file-card.selected').forEach(el => {
      el.classList.remove('selected');
    });
    selectedFile = null;
  }
});

// Right click on empty space in files grid
filesEl.addEventListener('contextmenu', (e) => {
  if (e.target === filesEl || e.target.classList.contains('empty-state')) {
    e.preventDefault();
    showContextMenu(e.clientX, e.clientY, null);
  }
});

// Context menu item handlers
document.getElementById('ctx-copy').onclick = () => {
  if (selectedFile) {
    copyToClipboard(selectedFile);
  }
  hideContextMenu();
};

document.getElementById('ctx-cut').onclick = () => {
  if (selectedFile) {
    cutToClipboard(selectedFile);
  }
  hideContextMenu();
};

document.getElementById('ctx-paste').onclick = () => {
  pasteFromClipboard();
  hideContextMenu();
};

document.getElementById('ctx-rename').onclick = () => {
  if (selectedFile) {
    showRenameModal(selectedFile);
  }
  hideContextMenu();
};

document.getElementById('ctx-delete').onclick = () => {
  if (selectedFile) {
    showDeleteModal(selectedFile);
  }
  hideContextMenu();
};

document.getElementById('ctx-download').onclick = () => {
  if (selectedFile && !selectedFile.is_dir) {
    downloadFile(selectedFile);
  }
  hideContextMenu();
};

// Rename modal handlers
document.getElementById('rename-cancel').onclick = hideRenameModal;
document.getElementById('rename-confirm').onclick = confirmRename;
document.getElementById('rename-input').onkeydown = (e) => {
  if (e.key === 'Enter') {
    confirmRename();
  } else if (e.key === 'Escape') {
    hideRenameModal();
  }
};

// Delete modal handlers
document.getElementById('delete-cancel').onclick = hideDeleteModal;
document.getElementById('delete-confirm').onclick = confirmDelete;

// Close modals on overlay click
renameModalEl.onclick = (e) => {
  if (e.target === renameModalEl) {
    hideRenameModal();
  }
};

deleteModalEl.onclick = (e) => {
  if (e.target === deleteModalEl) {
    hideDeleteModal();
  }
};

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
  // Don't handle shortcuts if modal is open or typing in input
  if (renameModalEl.classList.contains('visible') ||
    deleteModalEl.classList.contains('visible') ||
    e.target.tagName === 'INPUT') {
    return;
  }

  // Delete key
  if (e.key === 'Delete' && selectedFile) {
    e.preventDefault();
    showDeleteModal(selectedFile);
  }

  // F2 - Rename
  if (e.key === 'F2' && selectedFile) {
    e.preventDefault();
    showRenameModal(selectedFile);
  }

  // Ctrl+C - Copy
  if (e.ctrlKey && e.key === 'c' && selectedFile) {
    e.preventDefault();
    copyToClipboard(selectedFile);
  }

  // Ctrl+X - Cut
  if (e.ctrlKey && e.key === 'x' && selectedFile) {
    e.preventDefault();
    cutToClipboard(selectedFile);
  }

  // Ctrl+V - Paste
  if (e.ctrlKey && e.key === 'v' && clipboard.items.length > 0) {
    e.preventDefault();
    pasteFromClipboard();
  }

  // Escape - hide context menu
  if (e.key === 'Escape') {
    hideContextMenu();
  }
});

// ===== 13. Initial Load =====
loadExplorer();