import { displayGantt } from "./displayGantt.js";

async function changeFile(file) {
    const sheetSelector = document.getElementById('sheetSelector');
    sheetSelector.innerHTML = '';
    if (!file || !file.toLowerCase().endsWith('.xlsx')) {
        sheetSelector.style.display = 'none';
        return;
    }

    try {
        const resp = await fetch('/api/sheets?file=' + encodeURIComponent(file));
        if (!resp.ok) throw new Error('failed to fetch sheets');
        const data = await resp.json();
        const defaultOption = document.createElement('option');
        defaultOption.value = '';
        defaultOption.textContent = 'All Sheets (default)';
        sheetSelector.appendChild(defaultOption);
        data.sheets.forEach(s => {
            const o = document.createElement('option');
            o.value = s;
            o.textContent = s;
            sheetSelector.appendChild(o);
        });
        sheetSelector.style.display = '';
    } catch (err) {
        console.error('Error loading sheets:', err);
        sheetSelector.style.display = 'none';
    }
}

async function loadFiles() {
    try {
        const response = await fetch('/api/files');
        const data = await response.json();
        const selector = document.getElementById('fileSelector');
        data.files.forEach(file => {
            const option = document.createElement('option');
            option.value = file;
            option.textContent = file;
            selector.appendChild(option);
        });
        // When files loaded, also wire change event to fetch sheets for xlsx
        selector.addEventListener('change', () => changeFile(selector.value));
    } catch (error) {
        console.error('Error loading files:', error);
    }
}

async function loadRecords() {
    const selector = document.getElementById('fileSelector');
    const sheetSelector = document.getElementById('sheetSelector');
    const file = selector.value;
    const sheet = sheetSelector ? sheetSelector.value : '';
    let url = '/api/records';
    const params = [];
    if (file) params.push('file=' + encodeURIComponent(file));
    if (sheet) params.push('sheet=' + encodeURIComponent(sheet));
    if (params.length) url += '?' + params.join('&');
    try {
        const response = await fetch(url);
        if (!response.ok) throw new Error('failed to fetch records');
        const records = await response.json();
        return records;
    } catch (error) {
        console.error('Error loading data:', error);
        return null;
    }
}

function setupUserSelector(records) {
    const userSelector = document.getElementById('userSelector');
    userSelector.innerHTML = '';
    const users = Array.from(new Set(records.map(r => r.User))).sort();
    if (users.length === 0) {
        userSelector.style.display = 'none';
        return;
    }

    userSelector.style.display = '';
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'All Users (default)';
    userSelector.appendChild(defaultOption);
    users.forEach(user => {
        const o = document.createElement('option');
        o.value = user;
        o.textContent = user;
        userSelector.appendChild(o);
    });
}

let fullRecords = [];

async function getData() {
    const records = await loadRecords();
    if (!records) return;
    fullRecords = records;

    setupUserSelector(records);
    displayGantt(records);
}

function filterUser(e) {
    const user = e.target.value;
    const filteredRecords = user ? fullRecords.filter(r => r.User === user) : fullRecords;
    displayGantt(filteredRecords);
}

function downloadCsv() {
    const gantt = document.getElementById('gantt');
    if (!gantt) return;
    const table = gantt.querySelector('table');
    if (!table) return;

    const rows = Array.from(table.rows);
    const csvLines = rows.map(row => {
        const cells = Array.from(row.cells).map(cell => {
            let text = cell.textContent || '';
            text = text.replace(/\u00A0/g, ' ').trim();
            if (text.includes('"')) text = text.replace(/"/g, '""');
            if (text.includes(',') || text.includes('"') || text.includes('\n')) {
                return `"${text}"`;
            }
            return text;
        });
        return cells.join(',');
    });

    const csvContent = csvLines.join('\r\n');

    // ファイル名に選択中の file/sheet を反映（存在する場合）
    const selector = document.getElementById('fileSelector');
    const sheetSelector = document.getElementById('sheetSelector');
    const file = selector ? selector.value : '';
    const sheet = sheetSelector ? sheetSelector.value : '';
    const fnameParts = ['usage'];
    if (file) fnameParts.push(file.replace(/[^a-z0-9_\-\.]/gi, '_'));
    if (sheet) fnameParts.push(sheet.replace(/[^a-z0-9_\-\.]/gi, '_'));
    const filename = fnameParts.join('_') + '.csv';

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

document.addEventListener('DOMContentLoaded', loadFiles);
document.getElementById('fileSelector').addEventListener('change', getData);
document.getElementById('sheetSelector').addEventListener('change', getData);
document.getElementById('userSelector').addEventListener('change', filterUser);
document.getElementById('downloadCsv').addEventListener('click', downloadCsv);
