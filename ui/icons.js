let FILE_ICONS_MAP = { filenames: {}, extensions: {} };
let FOLDER_ICONS_MAP = {};

/**
 * Promise that resolves when icons are loaded
 */
const iconsLoaded = Promise.all([
  fetch('file_icons.json').then(res => res.json()),
  fetch('folder_icons.json').then(res => res.json())
]).then(([fileMap, folderMap]) => {
  FILE_ICONS_MAP = fileMap;
  FOLDER_ICONS_MAP = folderMap;
  console.log('Icon mappings loaded');
}).catch(err => {
  console.error('Failed to load icons:', err);
});

/**
 * Get the icon name for a file name
 * prioritize exact filename match, then extension match
 */
function getIconName(fileName) {
  if (!fileName) return 'file';
  const name = fileName.toLowerCase();

  // Check exact filename
  if (FILE_ICONS_MAP.filenames[name]) {
    return FILE_ICONS_MAP.filenames[name];
  }

  // Check extensions
  const parts = name.split('.');
  if (parts.length > 1) {
    for (let i = 1; i < parts.length; i++) {
      const ext = parts.slice(i).join('.');
      if (FILE_ICONS_MAP.extensions[ext]) {
        return FILE_ICONS_MAP.extensions[ext];
      }
    }
  }

  return 'file';
}

/**
 * Get the icon name for a folder name
 */
function getFolderIconName(folderName) {
  if (!folderName) return 'folder';
  const name = folderName.toLowerCase();

  if (FOLDER_ICONS_MAP[name]) {
    return FOLDER_ICONS_MAP[name];
  }

  return 'folder';
}
