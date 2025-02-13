import './style.css'
import { GetBingWallpapers, GetAndSetWallpaper, SetAutoStart } from '../wailsjs/go/main/App'

// 初始化应用
document.addEventListener('DOMContentLoaded', async () => {
    await loadWallpapers();

    document.querySelector('#refresh').addEventListener('click', async () => {
        await loadWallpapers();
    });
});

// 加载壁纸
async function loadWallpapers() {
    try {
        const wallpapers = await GetBingWallpapers();
        const container = document.getElementById('wallpapers');
        container.innerHTML = '';

        wallpapers.forEach(wallpaper => {
            const item = createWallpaperItem(wallpaper);
            container.appendChild(item);
        });
    } catch (error) {
        console.error('Failed to load wallpapers:', error);
    }
}

// 创建壁纸项
function createWallpaperItem(wallpaper) {
    const div = document.createElement('div');
    div.className = 'wallpaper-item';
    
    div.innerHTML = `
        <img src="${wallpaper.url}" alt="${wallpaper.title}" class="wallpaper-image">
        <div class="wallpaper-info">
            <div class="wallpaper-title">${wallpaper.title}</div>
            <div class="wallpaper-copyright">${wallpaper.copyright}</div>
            <div class="wallpaper-startdate">${wallpaper.startdate}</div>
        </div>
    `;

    div.querySelector('.wallpaper-image').addEventListener('click', async () => {
        try {
            await GetAndSetWallpaper(wallpaper.url);
            alert('壁纸设置成功！');
        } catch (error) {
            console.error('Failed to set wallpaper:', error);
            alert('设置壁纸失败：' + error);
        }
    });

    return div;
}