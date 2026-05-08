export function displayGantt(records) {
    // Sort records like CLI
    records.sort((a, b) => {
        if (a.EquipmentID === b.EquipmentID) {
            if (new Date(a.BeginDate) - new Date(b.BeginDate) === 0) {
                return new Date(a.EndDate) - new Date(b.EndDate);
            }
            return new Date(a.BeginDate) - new Date(b.BeginDate);
        }
        return a.EquipmentID.localeCompare(b.EquipmentID);
    });

    // Get min and max dates
    let minDate = new Date(records[0].BeginDate);
    let maxDate = new Date(records[0].EndDate);
    records.forEach(record => {
        const begin = new Date(record.BeginDate);
        const end = new Date(record.EndDate);
        if (begin < minDate) minDate = begin;
        if (end > maxDate) maxDate = end;
    });

    // Make labels
    const dateLabels = [];
    const weekLabels = [];
    for (let d = new Date(minDate); d <= maxDate; d.setDate(d.getDate() + 1)) {
        dateLabels.push(d.getDate().toString().padStart(2, '0'));
        const weekday = d.getDay();
        const weekdays = ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'];
        weekLabels.push(weekdays[weekday]);
    }

    // Build table
    const gantt = document.getElementById('gantt');
    gantt.innerHTML = '';

    const table = document.createElement('table');

    // Data rows
    let prevID = '';
    records.forEach(record => {
        if (record.EquipmentID !== prevID) {
            createHeaderRow(table, dateLabels, weekLabels);
            createWeekRow(table, weekLabels);
        }
        const row = table.insertRow();
        const userCell = row.insertCell();
        userCell.textContent = record.User;
        userCell.style.width = '100px';

        let idx = 0;
        const isEndless = record.EndDate === '0001-01-01T00:00:00Z' || !record.EndDate;
        for (let d = new Date(minDate); d <= maxDate; d.setDate(d.getDate() + 1)) {
            const cell = row.insertCell();
            const w = weekLabels[idx];
            if (d >= new Date(record.BeginDate) && (isEndless || d <= new Date(record.EndDate))) {
                cell.textContent = isEndless ? '??' : '**';
                cell.className = 'usage ' + w.toLowerCase();
            } else {
                cell.textContent = '  ';
                cell.className = w.toLowerCase();
            }
            idx++;
        }
        prevID = record.EquipmentID;
    });

    gantt.appendChild(table);
}
// Helper: create header row.
function createHeaderRow(table, dateLabels, weekLabels) {
    const row = table.insertRow();
    const nameCell = row.insertCell();
    nameCell.textContent = 'User Name';
    nameCell.style.fontWeight = 'bold';
    dateLabels.forEach((label, i) => {
        const cell = row.insertCell();
        cell.textContent = label;
        cell.className = 'date-label ' + weekLabels[i].toLowerCase();
    });
}
// Helper: create week row.
function createWeekRow(table, weekLabels) {
    const row = table.insertRow();
    row.insertCell().textContent = '';
    weekLabels.forEach(w => {
        const cell = row.insertCell();
        cell.textContent = w;
        cell.className = 'week-label ' + w.toLowerCase();
    });
}
