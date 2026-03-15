// ─── HabitFlow SPA ──────────────────────────────────────────────
(function () {
  'use strict';

  const $ = (s, p) => (p || document).querySelector(s);
  const $$ = (s, p) => [...(p || document).querySelectorAll(s)];
  const app = () => $('#app');

  // ── Toast ──
  function toast(msg, type = 'success') {
    const c = $('#toast-container');
    const colors = { success: 'bg-primary', error: 'bg-red-500', info: 'bg-slate-700' };
    const el = document.createElement('div');
    el.className = `toast-in ${colors[type] || colors.info} text-white px-4 py-2.5 rounded-xl shadow-lg text-sm font-medium max-w-xs text-center`;
    el.textContent = msg;
    c.appendChild(el);
    setTimeout(() => { el.classList.replace('toast-in', 'toast-out'); setTimeout(() => el.remove(), 300); }, 2500);
  }

  function escapeHtml(str) { const d = document.createElement('div'); d.textContent = str; return d.innerHTML; }

  // ── Date helpers ──
  const dayNames = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
  const monthNames = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'];

  function formatDateID(d) {
    const dt = d ? new Date(d) : new Date();
    return `${dayNames[dt.getDay()]}, ${dt.getDate()} ${monthNames[dt.getMonth()]} ${dt.getFullYear()}`;
  }

  function formatDateShort(dateStr) {
    if (!dateStr) return '';
    const parts = dateStr.split('-');
    const d = new Date(parts[0], parts[1] - 1, parts[2]);
    return `${d.getDate()} ${monthNames[d.getMonth()].slice(0, 3)}`;
  }

  function dateStr(d) {
    return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-' + String(d.getDate()).padStart(2, '0');
  }

  // ── Category helpers ──
  function categoryColor(cat) {
    const map = {
      health: { bg: 'bg-red-900/30', text: 'text-red-400', light: 'bg-red-100 text-red-500' },
      fitness: { bg: 'bg-orange-900/30', text: 'text-orange-400', light: 'bg-orange-100 text-orange-500' },
      learning: { bg: 'bg-purple-900/30', text: 'text-purple-400', light: 'bg-purple-100 text-purple-500' },
      productivity: { bg: 'bg-amber-900/30', text: 'text-amber-400', light: 'bg-amber-100 text-amber-500' },
      mindfulness: { bg: 'bg-teal-900/30', text: 'text-teal-400', light: 'bg-teal-100 text-teal-500' },
      finance: { bg: 'bg-green-900/30', text: 'text-green-400', light: 'bg-green-100 text-green-500' },
      social: { bg: 'bg-pink-900/30', text: 'text-pink-400', light: 'bg-pink-100 text-pink-500' },
      general: { bg: 'bg-slate-800', text: 'text-slate-400', light: 'bg-slate-200 text-slate-500' },
      spiritual: { bg: 'bg-teal-900/30', text: 'text-teal-400', light: 'bg-teal-100 text-teal-500' },
      creative: { bg: 'bg-rose-900/30', text: 'text-rose-400', light: 'bg-rose-100 text-rose-500' },
    };
    return map[cat] || map.general;
  }
  function categoryLabel(cat) {
    const map = { health: 'Kesehatan', fitness: 'Olahraga', learning: 'Belajar', productivity: 'Produktivitas', mindfulness: 'Mindfulness', finance: 'Keuangan', social: 'Sosial', general: 'Lainnya', spiritual: 'Mindfulness', creative: 'Kreatif' };
    return map[cat] || 'Lainnya';
  }
  function categorySVG(cat) {
    const svgs = {
      health: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/></svg>',
      fitness: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.4 14.4 9.6 9.6"/><path d="M18.657 21.485a2 2 0 1 1-2.829-2.828l-1.767 1.768a2 2 0 1 1-2.829-2.829l6.364-6.364a2 2 0 1 1 2.829 2.829l-1.768 1.767a2 2 0 1 1 2.828 2.829z"/><path d="m21.5 21.5-1.4-1.4"/><path d="M3.9 3.9 2.5 2.5"/><path d="M6.404 12.768a2 2 0 1 1-2.829-2.829l1.768-1.767a2 2 0 1 1-2.828-2.829l2.828-2.828a2 2 0 1 1 2.829 2.828l1.767-1.768a2 2 0 1 1 2.829 2.829z"/></svg>',
      learning: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>',
      productivity: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="6"/><circle cx="12" cy="12" r="2"/></svg>',
      mindfulness: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5a3 3 0 1 0-5.997.125 4 4 0 0 0-2.526 5.77 4 4 0 0 0 .556 6.588A4 4 0 1 0 12 18Z"/><path d="M12 5a3 3 0 1 1 5.997.125 4 4 0 0 1 2.526 5.77 4 4 0 0 1-.556 6.588A4 4 0 1 1 12 18Z"/><path d="M15 13a4.5 4.5 0 0 1-3-4 4.5 4.5 0 0 1-3 4"/><path d="M12 18v-5.5"/><path d="M9.5 9 11 14.5"/><path d="M14.5 9 13 14.5"/></svg>',
      finance: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12V7H5a2 2 0 0 1 0-4h14v4"/><path d="M3 5v14a2 2 0 0 0 2 2h16v-5"/><path d="M18 12a2 2 0 0 0 0 4h4v-4Z"/></svg>',
      social: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>',
      general: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="7" x="3" y="3" rx="1"/><rect width="7" height="7" x="14" y="3" rx="1"/><rect width="7" height="7" x="14" y="14" rx="1"/><rect width="7" height="7" x="3" y="14" rx="1"/></svg>',
    };
    if (cat === 'spiritual') return svgs.mindfulness;
    if (cat === 'creative') return svgs.general;
    return svgs[cat] || svgs.general;
  }
  const chevronDownSVG = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>';

  function initCategoryDropdown() {
    const trigger = $('#category-trigger');
    const panel = $('#category-panel');
    const hiddenInput = $('#f-category');
    const display = $('#category-display');
    const chevron = $('#category-chevron');
    if (!trigger || !panel) return;
    let isOpen = false;
    const options = $$('.category-option', panel);
    let focusIdx = -1;

    function open() {
      isOpen = true;
      panel.classList.remove('hidden');
      panel.classList.add('dropdown-enter');
      chevron.classList.add('rotate-180');
      trigger.setAttribute('aria-expanded', 'true');
      focusIdx = -1;
    }
    function close() {
      isOpen = false;
      panel.classList.add('dropdown-exit');
      chevron.classList.remove('rotate-180');
      trigger.setAttribute('aria-expanded', 'false');
      setTimeout(() => {
        panel.classList.add('hidden');
        panel.classList.remove('dropdown-enter', 'dropdown-exit');
      }, 150);
      focusIdx = -1;
    }
    function selectOption(opt) {
      const value = opt.dataset.value;
      const svg = categorySVG(value);
      const text = categoryLabel(value);
      display.innerHTML = svg + '<span>' + text + '</span>';
      hiddenInput.value = value;
      options.forEach(o => o.classList.remove('bg-primary/10', 'text-primary'));
      opt.classList.add('bg-primary/10', 'text-primary');
      close();
      trigger.focus();
    }
    function focusOption(idx) {
      if (idx < 0) idx = options.length - 1;
      if (idx >= options.length) idx = 0;
      focusIdx = idx;
      options.forEach(o => o.classList.remove('bg-primary/10'));
      options[focusIdx].classList.add('bg-primary/10');
      options[focusIdx].scrollIntoView({ block: 'nearest' });
    }

    trigger.addEventListener('click', (e) => { e.preventDefault(); isOpen ? close() : open(); });
    options.forEach(opt => {
      opt.addEventListener('click', (e) => { e.preventDefault(); selectOption(opt); });
    });
    document.addEventListener('click', (e) => {
      if (isOpen && !trigger.contains(e.target) && !panel.contains(e.target)) close();
    });
    trigger.addEventListener('keydown', (e) => {
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp') { e.preventDefault(); if (!isOpen) open(); focusOption(e.key === 'ArrowDown' ? 0 : options.length - 1); }
      if (e.key === 'Escape' && isOpen) { e.preventDefault(); close(); }
      if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); isOpen ? close() : open(); }
    });
    panel.addEventListener('keydown', (e) => {
      if (e.key === 'ArrowDown') { e.preventDefault(); focusOption(focusIdx + 1); }
      else if (e.key === 'ArrowUp') { e.preventDefault(); focusOption(focusIdx - 1); }
      else if ((e.key === 'Enter' || e.key === ' ') && focusIdx >= 0) { e.preventDefault(); selectOption(options[focusIdx]); }
      else if (e.key === 'Escape') { e.preventDefault(); close(); trigger.focus(); }
      else if (e.key === 'Tab') { close(); }
    });
    // Mark initial selected option
    const currentVal = hiddenInput.value;
    options.forEach(o => { if (o.dataset.value === currentVal) o.classList.add('bg-primary/10', 'text-primary'); });
  }

  // ── Time Dropdown helpers ──
  const clockSVG = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>';
  const timeSVGs = {
    sunrise: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v8"/><path d="m4.93 10.93 1.41 1.41"/><path d="M2 18h2"/><path d="M20 18h2"/><path d="m19.07 10.93-1.41 1.41"/><path d="M22 22H2"/><path d="m8 6 4-4 4 4"/><path d="M16 18a4 4 0 0 0-8 0"/></svg>',
    sun: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="m4.93 4.93 1.41 1.41"/><path d="m17.66 17.66 1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="m6.34 17.66-1.41 1.41"/><path d="m19.07 4.93-1.41 1.41"/></svg>',
    sunset: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 10V2"/><path d="m4.93 10.93 1.41 1.41"/><path d="M2 18h2"/><path d="M20 18h2"/><path d="m19.07 10.93-1.41 1.41"/><path d="M22 22H2"/><path d="m8 6 4 4 4-4"/><path d="M16 18a4 4 0 0 0-8 0"/></svg>',
    moon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/></svg>',
    none: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="m4.9 4.9 14.2 14.2"/></svg>',
  };
  function timeIcon(hour) {
    if (hour < 0) return timeSVGs.none;
    if (hour < 10) return timeSVGs.sunrise;
    if (hour < 15) return timeSVGs.sun;
    if (hour < 18) return timeSVGs.sunset;
    return timeSVGs.moon;
  }
  function timeLabel(hour) {
    if (hour < 10) return 'Pagi';
    if (hour < 15) return 'Siang';
    if (hour < 18) return 'Sore';
    return 'Malam';
  }
  const notifyTimes = [
    { value: '', label: 'Tidak ada notifikasi', hour: -1 },
    { value: '05:00', label: '05:00', hour: 5 },
    { value: '06:00', label: '06:00', hour: 6 },
    { value: '07:00', label: '07:00', hour: 7 },
    { value: '08:00', label: '08:00', hour: 8 },
    { value: '09:00', label: '09:00', hour: 9 },
    { value: '10:00', label: '10:00', hour: 10 },
    { value: '11:00', label: '11:00', hour: 11 },
    { value: '12:00', label: '12:00', hour: 12 },
    { value: '13:00', label: '13:00', hour: 13 },
    { value: '14:00', label: '14:00', hour: 14 },
    { value: '15:00', label: '15:00', hour: 15 },
    { value: '16:00', label: '16:00', hour: 16 },
    { value: '17:00', label: '17:00', hour: 17 },
    { value: '18:00', label: '18:00', hour: 18 },
    { value: '19:00', label: '19:00', hour: 19 },
    { value: '20:00', label: '20:00', hour: 20 },
    { value: '21:00', label: '21:00', hour: 21 },
    { value: '22:00', label: '22:00', hour: 22 },
  ];

  function isPresetNotifyTime(val) {
    return notifyTimes.some(t => t.value === (val || ''));
  }

  function notifyTimeDisplayHTML(val) {
    if (!val) return timeSVGs.none + '<span class="text-slate-400">Tidak ada notifikasi</span>';
    const t = notifyTimes.find(t => t.value === val);
    const hour = t ? t.hour : parseInt(val.split(':')[0], 10);
    return timeIcon(hour) + '<span>' + val + ' — ' + timeLabel(hour) + '</span>';
  }
  function buildTimeOptions(selected) {
    let prevGroup = '';
    return notifyTimes.map(t => {
      let separator = '';
      let group = '';
      if (t.hour < 0) group = '';
      else if (t.hour < 10) group = 'Pagi';
      else if (t.hour < 15) group = 'Siang';
      else if (t.hour < 18) group = 'Sore';
      else group = 'Malam';
      if (group && group !== prevGroup) {
        separator = `<div class="px-4 pt-3 pb-1 text-[11px] font-bold uppercase tracking-widest text-slate-500">${group}</div>`;
        prevGroup = group;
      }
      const isSelected = t.value === (selected || '');
      return separator + `
        <button type="button" class="time-option mx-2 my-0.5 w-[calc(100%-1rem)] rounded-lg border border-transparent
                                      flex items-center justify-between gap-3 px-3 py-2.5
                                      text-sm text-slate-200 hover:bg-primary/10 hover:border-primary/20
                                      focus:outline-none transition-colors
                                      ${isSelected ? 'bg-primary/10 text-primary border-primary/30 time-option-selected' : ''}"
                data-value="${t.value}" role="option">
          <span class="flex items-center gap-3">
            ${timeIcon(t.hour)}
            <span>${t.hour < 0 ? 'Tidak ada notifikasi' : t.label + ' — ' + timeLabel(t.hour)}</span>
          </span>
          <span class="time-option-check ${isSelected ? 'opacity-100' : 'opacity-0'}">${Icons.check(16, 'currentColor')}</span>
        </button>`;
    }).join('');
  }

  function customTimeDisplayHTML(val) {
    if (!val) return clockSVG + '<span class="text-slate-400">Pilih waktu sendiri</span>';
    const hour = parseInt(val.split(':')[0], 10);
    return timeIcon(hour) + '<span>' + val + ' — ' + timeLabel(hour) + '</span>';
  }

  function buildCustomTimeOptions(selected) {
    const selectedValue = selected || '';
    const manualBadge = selectedValue
      ? `<span class="text-[11px] px-2 py-1 rounded-md bg-primary/15 text-primary font-semibold">Dipilih: ${selectedValue}</span>`
      : '<span class="text-[11px] text-slate-500">Masukkan jam dan menit bebas</span>';

    return `
      <div class="px-3 pt-2 pb-2 border-b border-slate-800/80">
        <div class="rounded-xl border border-slate-700/90 bg-slate-800/40 px-3 py-3">
          <div class="flex items-center justify-between mb-2">
            <p class="text-xs font-semibold tracking-wide text-slate-300">Waktu Bebas</p>
            ${manualBadge}
          </div>
          <div class="flex items-center gap-2">
            <input id="custom-hour-input" type="text" inputmode="numeric" maxlength="2" placeholder="09"
              class="custom-time-field w-14 h-10 rounded-lg border border-slate-700 bg-slate-900/70 text-center text-sm font-semibold text-slate-100 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary">
            <span class="text-slate-400 font-bold">:</span>
            <input id="custom-minute-input" type="text" inputmode="numeric" maxlength="2" placeholder="30"
              class="custom-time-field w-14 h-10 rounded-lg border border-slate-700 bg-slate-900/70 text-center text-sm font-semibold text-slate-100 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary">
            <button type="button" id="custom-time-apply"
              class="ml-auto h-10 px-3 rounded-lg bg-primary/90 text-white text-xs font-semibold hover:bg-primary transition-colors">
              Terapkan
            </button>
          </div>
        </div>
      </div>
    `;
  }

  function initTimeDropdown() {
    const trigger = $('#time-trigger');
    const panel = $('#time-panel');
    const hiddenInput = $('#f-notify');
    const display = $('#time-display');
    const chevron = $('#time-chevron');
    if (!trigger || !panel) return;
    let isOpen = false;
    const options = $$('.time-option', panel);
    let focusIdx = -1;

    function open() {
      isOpen = true;
      panel.classList.remove('hidden');
      panel.classList.add('dropdown-enter');
      chevron.classList.add('rotate-180');
      trigger.setAttribute('aria-expanded', 'true');
      focusIdx = -1;
      // Scroll to selected option
      const sel = panel.querySelector('.bg-primary\\/10');
      if (sel) setTimeout(() => sel.scrollIntoView({ block: 'nearest' }), 50);
    }
    function close() {
      isOpen = false;
      panel.classList.add('dropdown-exit');
      chevron.classList.remove('rotate-180');
      trigger.setAttribute('aria-expanded', 'false');
      setTimeout(() => {
        panel.classList.add('hidden');
        panel.classList.remove('dropdown-enter', 'dropdown-exit');
      }, 150);
      focusIdx = -1;
    }
    function selectOption(opt) {
      const value = opt.dataset.value;
      display.innerHTML = notifyTimeDisplayHTML(value);
      hiddenInput.value = value;
      hiddenInput.dispatchEvent(new Event('change', { bubbles: true }));
      options.forEach(o => {
        o.classList.remove('bg-primary/10', 'text-primary', 'border-primary/30', 'time-option-selected');
        const check = o.querySelector('.time-option-check');
        if (check) check.classList.remove('opacity-100');
        if (check) check.classList.add('opacity-0');
      });
      opt.classList.add('bg-primary/10', 'text-primary', 'border-primary/30', 'time-option-selected');
      const selectedCheck = opt.querySelector('.time-option-check');
      if (selectedCheck) {
        selectedCheck.classList.remove('opacity-0');
        selectedCheck.classList.add('opacity-100');
      }
      close();
      trigger.focus();
    }
    function focusOption(idx) {
      if (idx < 0) idx = options.length - 1;
      if (idx >= options.length) idx = 0;
      focusIdx = idx;
      options.forEach(o => o.classList.remove('time-option-focus'));
      options[focusIdx].classList.add('time-option-focus');
      options[focusIdx].scrollIntoView({ block: 'nearest' });
    }

    trigger.addEventListener('click', (e) => { e.preventDefault(); isOpen ? close() : open(); });
    options.forEach(opt => {
      opt.addEventListener('click', (e) => { e.preventDefault(); selectOption(opt); });
    });
    document.addEventListener('click', (e) => {
      if (isOpen && !trigger.contains(e.target) && !panel.contains(e.target)) close();
    });
    trigger.addEventListener('keydown', (e) => {
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp') { e.preventDefault(); if (!isOpen) open(); focusOption(e.key === 'ArrowDown' ? 0 : options.length - 1); }
      if (e.key === 'Escape' && isOpen) { e.preventDefault(); close(); }
      if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); isOpen ? close() : open(); }
    });
    panel.addEventListener('keydown', (e) => {
      if (e.key === 'ArrowDown') { e.preventDefault(); focusOption(focusIdx + 1); }
      else if (e.key === 'ArrowUp') { e.preventDefault(); focusOption(focusIdx - 1); }
      else if ((e.key === 'Enter' || e.key === ' ') && focusIdx >= 0) { e.preventDefault(); selectOption(options[focusIdx]); }
      else if (e.key === 'Escape') { e.preventDefault(); close(); trigger.focus(); }
      else if (e.key === 'Tab') { close(); }
    });
  }

  function skeleton(n) { return Array(n).fill('<div class="skeleton h-20 w-full mb-3"></div>').join(''); }

  // ── SVG Helpers ──
  const googleSVG = `<svg class="h-5 w-5" viewBox="0 0 24 24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>`;

  // ════════════════════════════════════════════════════════════════
  // ROUTER
  // ════════════════════════════════════════════════════════════════
  function route() {
    const hash = location.hash || (API.isLoggedIn() ? '#/dashboard' : '#/login');
    const path = hash.slice(1);
    if (path === '/login') return renderLogin();
    if (path === '/register') return renderRegister();
    if (!API.isLoggedIn() && path !== '/login' && path !== '/register') { location.hash = '#/login'; return; }
    if (API.isLoggedIn() && (path === '/login' || path === '/register')) { location.hash = '#/dashboard'; return; }
    if (path === '/onboarding') return renderOnboarding();
    if (path === '/dashboard' || path === '/') return renderDashboard();
    if (path === '/habits') return renderHabits();
    if (path === '/habits/new') return renderHabitForm();
    if (path.match(/^\/habits\/\d+\/edit$/)) return renderHabitForm(path.match(/\d+/)[0]);
    if (path === '/report') return renderReport();
    if (path === '/settings') return renderSettings();
    return renderDashboard();
  }
  window.addEventListener('hashchange', route);
  window.addEventListener('DOMContentLoaded', route);

  // ════════════════════════════════════════════════════════════════
  // LAYOUT COMPONENTS
  // ════════════════════════════════════════════════════════════════
  function sidebarNav(active) {
    const items = [
      { id: 'dashboard', svg: Icons.calendarDays, label: 'Hari Ini' },
      { id: 'habits', svg: Icons.checkCircle, label: 'Habit Saya' },
      { id: 'report', svg: Icons.barChart3, label: 'Laporan' },
      { id: 'settings', svg: Icons.settings, label: 'Pengaturan' },
    ];
    return items.map(i => {
      const a = active === i.id;
      return `<a href="#/${i.id}" class="flex items-center gap-3 px-4 py-3 rounded-xl ${a ? 'bg-primary/10 text-primary font-semibold' : 'text-slate-400 hover:bg-slate-800 transition-colors'}">
        ${i.svg(20)}<span class="text-sm">${i.label}</span></a>`;
    }).join('');
  }

  function mobileNav(active) {
    const items = [
      { id: 'dashboard', svg: Icons.calendarDays, label: 'Hari Ini' },
      { id: 'habits', svg: Icons.checkCircle, label: 'Habit' },
      { id: 'report', svg: Icons.barChart3, label: 'Laporan' },
      { id: 'settings', svg: Icons.user, label: 'Profil' },
    ];
    return items.map(i => {
      const a = active === i.id;
      return `<a href="#/${i.id}" class="flex flex-col items-center gap-1 ${a ? 'text-primary' : 'text-slate-400'}">
        ${i.svg(20)}
        <span class="text-[10px] ${a ? 'font-bold' : 'font-medium'}">${i.label}</span></a>`;
    }).join('');
  }

  function appShell(activePage, headerHTML, contentHTML, opts = {}) {
    const user = API.getUser();
    const initials = (user?.name || 'U').split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2);
    const showFab = opts.fab !== false && activePage === 'dashboard';
    return `
    <div class="flex h-screen overflow-hidden">
      <aside class="hidden md:flex w-64 flex-col bg-slate-900 border-r border-slate-800 shrink-0">
        <div class="p-6 flex items-center gap-3">
          <div class="size-8 bg-primary rounded-lg flex items-center justify-center text-white">
            ${Icons.repeat(20, '#fff')}
          </div>
          <h2 class="text-xl font-bold tracking-tight text-primary">HabitFlow</h2>
        </div>
        <nav class="flex-1 px-4 space-y-1 mt-2">${sidebarNav(activePage)}</nav>
        <div class="p-4 mt-auto">
          <div class="bg-slate-800 rounded-xl p-4 flex items-center gap-3">
            <div class="size-10 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-sm">${initials}</div>
            <div class="overflow-hidden">
              <p class="text-sm font-bold truncate">${escapeHtml(user?.name || 'User')}</p>
              <p class="text-xs text-slate-500 truncate">${escapeHtml(user?.email || '')}</p>
            </div>
          </div>
        </div>
      </aside>
      <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <header class="bg-background-dark/80 backdrop-blur-md sticky top-0 z-10 px-6 py-4 border-b border-slate-800 flex items-center justify-between">${headerHTML}</header>
        <div class="flex-1 overflow-y-auto p-6 space-y-6 pb-24 md:pb-6 hide-scrollbar">${contentHTML}</div>
      </main>
      <nav class="md:hidden bg-slate-900 border-t border-slate-800 px-6 py-3 flex justify-around items-center fixed bottom-0 left-0 right-0 z-30 safe-area-bottom">${mobileNav(activePage)}</nav>
      ${showFab ? `<a href="#/habits/new" class="fixed bottom-24 right-6 md:bottom-8 md:right-8 size-14 bg-primary text-white rounded-full shadow-lg shadow-primary/40 flex items-center justify-center hover:scale-105 active:scale-95 transition-transform z-20 tap-target">${Icons.plus(28, '#fff')}</a>` : ''}
    </div>`;
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: LOGIN
  // ════════════════════════════════════════════════════════════════
  function renderLogin() {
    app().innerHTML = `
    <div class="relative flex min-h-screen w-full flex-col items-center justify-center overflow-x-hidden p-4">
      <div class="fixed inset-0 z-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-[10%] -left-[10%] w-[50%] h-[50%] rounded-full bg-primary/10 blur-[120px]"></div>
        <div class="absolute -bottom-[10%] -right-[10%] w-[50%] h-[50%] rounded-full bg-primary/5 blur-[120px]"></div>
      </div>
      <div class="relative z-10 w-full max-w-[440px] flex flex-col items-center">
        <div class="mb-8 flex flex-col items-center gap-2">
          <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-primary text-white shadow-lg shadow-primary/20">
            ${Icons.repeat(32, '#fff')}
          </div>
          <h1 class="text-2xl font-bold tracking-tight text-white">HabitFlow</h1>
        </div>
        <div class="w-full rounded-2xl bg-slate-900/50 backdrop-blur-sm p-8 shadow-xl border border-slate-800/50">
          <div class="mb-8">
            <h2 class="text-2xl font-bold leading-tight text-white">Masuk ke HabitFlow</h2>
            <p class="mt-2 text-slate-400">Silakan masukkan detail akun Anda</p>
          </div>
          <form id="auth-form" class="space-y-5">
            <div class="flex flex-col gap-2">
              <label class="text-sm font-medium text-slate-300">Email</label>
              <div class="group relative flex w-full items-center">
                <div class="absolute left-3 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.mail(20)}</div>
                <input id="f-email" type="email" placeholder="nama@email.com" required class="flex h-12 w-full rounded-xl border border-slate-700 bg-slate-800/50 px-10 py-2 text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-slate-100">
              </div>
            </div>
            <div class="flex flex-col gap-2">
              <div class="flex items-center justify-between">
                <label class="text-sm font-medium text-slate-300">Kata Sandi</label>
                <a class="text-xs font-medium text-primary hover:underline cursor-pointer">Lupa sandi?</a>
              </div>
              <div class="group relative flex w-full items-center">
                <div class="absolute left-3 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.lock(20)}</div>
                <input id="f-password" type="password" placeholder="••••••••" required minlength="8" class="flex h-12 w-full rounded-xl border border-slate-700 bg-slate-800/50 px-10 py-2 text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-slate-100">
                <button type="button" id="btn-toggle-pw" class="absolute right-3 text-slate-400 hover:text-slate-200">${Icons.eye(20)}</button>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <input class="h-4 w-4 rounded border-slate-700 bg-transparent text-primary focus:ring-primary" id="remember" type="checkbox">
              <label class="text-sm font-medium text-slate-400 select-none" for="remember">Ingat saya di perangkat ini</label>
            </div>
            <button type="submit" id="btn-submit" class="flex w-full items-center justify-center rounded-xl bg-primary h-12 px-4 text-sm font-bold text-white transition-all hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary/50 active:scale-[0.98] shadow-lg shadow-primary/20">Masuk</button>
          </form>
          <p id="auth-error" class="text-red-400 text-xs text-center mt-3 hidden"></p>
          <div class="relative my-8">
            <div class="absolute inset-0 flex items-center"><span class="w-full border-t border-slate-800"></span></div>
            <div class="relative flex justify-center text-xs uppercase"><span class="bg-slate-900 px-2 text-slate-500">Atau masuk dengan</span></div>
          </div>
          <button class="flex w-full items-center justify-center gap-2 rounded-xl border border-slate-700 bg-slate-800/50 h-11 px-4 text-sm font-medium text-slate-200 hover:bg-slate-700/50 transition-colors">${googleSVG} Google</button>
        </div>
        <p class="mt-8 text-center text-sm text-slate-400">Belum punya akun? <a href="#/register" class="font-bold text-primary hover:underline">Daftar sekarang</a></p>
        <div class="mt-12 flex gap-6 text-xs text-slate-400">
          <a class="hover:text-primary" href="#">Ketentuan Layanan</a>
          <a class="hover:text-primary" href="#">Kebijakan Privasi</a>
          <a class="hover:text-primary" href="#">Bantuan</a>
        </div>
      </div>
    </div>`;
    $('#btn-toggle-pw').onclick = () => {
      const pw = $('#f-password'), btn = $('#btn-toggle-pw');
      if (pw.type === 'password') { pw.type = 'text'; btn.innerHTML = Icons.eyeOff(20); }
      else { pw.type = 'password'; btn.innerHTML = Icons.eye(20); }
    };
    $('#auth-form').onsubmit = async (e) => {
      e.preventDefault();
      const btn = $('#btn-submit'), errEl = $('#auth-error');
      btn.disabled = true; btn.textContent = 'Memproses...'; errEl.classList.add('hidden');
      try {
        await API.login($('#f-email').value.trim(), $('#f-password').value);
        toast('Selamat datang!'); location.hash = '#/dashboard';
      } catch (err) { errEl.textContent = err.message; errEl.classList.remove('hidden'); btn.disabled = false; btn.textContent = 'Masuk'; }
    };
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: REGISTER
  // ════════════════════════════════════════════════════════════════
  function renderRegister() {
    app().innerHTML = `
    <div class="min-h-screen flex flex-col bg-background-dark">
      <header class="w-full border-b border-slate-800 bg-slate-900/80 backdrop-blur-sm fixed top-0 z-50">
        <div class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <div class="text-primary">${Icons.repeat(28, '#2ec2b3')}</div>
            <h2 class="text-xl font-bold tracking-tight">HabitFlow</h2>
          </div>
          <a href="#/login" class="inline-flex h-9 items-center justify-center rounded-lg bg-primary/10 px-4 text-sm font-bold text-primary hover:bg-primary/20 transition-colors">Masuk</a>
        </div>
      </header>
      <main class="flex-1 flex items-center justify-center pt-24 pb-12 px-4">
        <div class="w-full max-w-[520px] bg-slate-900 rounded-xl shadow-xl border border-slate-800 overflow-hidden">
          <div class="relative h-32 bg-primary/10 overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-br from-primary/20 to-transparent"></div>
            <div class="absolute -right-4 -top-4 size-32 bg-primary/10 rounded-full blur-2xl"></div>
            <div class="absolute left-8 bottom-4">
              <h1 class="text-3xl font-black tracking-tight text-white">Buat Akun Baru</h1>
              <p class="text-slate-400 text-sm mt-1">Mulai bangun kebiasaan positif dalam 30 detik.</p>
            </div>
          </div>
          <form id="reg-form" class="p-8 space-y-5">
            <div class="space-y-1.5">
              <label class="text-sm font-semibold text-slate-300 ml-1">Nama Lengkap</label>
              <div class="relative group">
                <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.user(20)}</span>
                <input id="r-name" type="text" required minlength="2" placeholder="Masukkan nama lengkap Anda" class="w-full pl-11 pr-4 py-3.5 rounded-lg border border-slate-800 bg-slate-900/50 focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-slate-500 text-white">
              </div>
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-semibold text-slate-300 ml-1">Email</label>
              <div class="relative group">
                <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.mail(20)}</span>
                <input id="r-email" type="email" required placeholder="nama@email.com" class="w-full pl-11 pr-4 py-3.5 rounded-lg border border-slate-800 bg-slate-900/50 focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-slate-500 text-white">
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-sm font-semibold text-slate-300 ml-1">Kata Sandi</label>
                <div class="relative group">
                  <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.lock(20)}</span>
                  <input id="r-password" type="password" required minlength="8" placeholder="••••••••" class="w-full pl-11 pr-4 py-3.5 rounded-lg border border-slate-800 bg-slate-900/50 focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-slate-500 text-white">
                </div>
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-semibold text-slate-300 ml-1">Konfirmasi Sandi</label>
                <div class="relative group">
                  <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.lockKeyhole(20)}</span>
                  <input id="r-confirm" type="password" required minlength="8" placeholder="••••••••" class="w-full pl-11 pr-4 py-3.5 rounded-lg border border-slate-800 bg-slate-900/50 focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-slate-500 text-white">
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2 px-1">
              <input class="rounded border-slate-700 text-primary focus:ring-primary cursor-pointer" id="terms" type="checkbox">
              <label class="text-xs text-slate-400 leading-tight" for="terms">Saya setuju dengan <a class="text-primary hover:underline" href="#">Syarat &amp; Ketentuan</a> serta <a class="text-primary hover:underline" href="#">Kebijakan Privasi</a> HabitFlow.</label>
            </div>
            <button type="submit" id="btn-reg" class="w-full bg-primary hover:bg-primary/90 text-slate-900 font-bold py-4 rounded-xl shadow-lg shadow-primary/20 transition-all flex items-center justify-center gap-2 active:scale-[0.98]">
              Daftar Sekarang ${Icons.arrowRight(20, '#fff')}
            </button>
            <p id="reg-error" class="text-red-400 text-xs text-center hidden"></p>
            <div class="pt-4 text-center"><p class="text-slate-400 text-sm">Sudah punya akun? <a href="#/login" class="text-primary font-bold hover:underline ml-1">Masuk</a></p></div>
          </form>
          <div class="px-8 pb-4 flex items-center justify-center gap-4">
            <div class="h-px bg-slate-800 flex-1"></div>
            <span class="text-[10px] uppercase tracking-widest text-slate-400 font-bold">Atau daftar dengan</span>
            <div class="h-px bg-slate-800 flex-1"></div>
          </div>
          <div class="px-8 pb-8 grid grid-cols-2 gap-4">
            <button class="flex items-center justify-center gap-2 py-2.5 border border-slate-800 rounded-lg hover:bg-slate-800/80 transition-colors text-sm font-medium bg-slate-900/50">${googleSVG} Google</button>
            <button class="flex items-center justify-center gap-2 py-2.5 border border-slate-800 rounded-lg hover:bg-slate-800/80 transition-colors text-sm font-medium bg-slate-900/50">
              ${Icons.linkedin(20, '#3b82f6')} LinkedIn
            </button>
          </div>
        </div>
      </main>
      <footer class="py-6 text-center text-slate-600 text-xs"><p>&copy; 2026 HabitFlow. Seluruh hak cipta dilindungi.</p></footer>
    </div>`;
    $('#reg-form').onsubmit = async (e) => {
      e.preventDefault();
      const errEl = $('#reg-error'), btn = $('#btn-reg');
      errEl.classList.add('hidden');
      const pw = $('#r-password').value, confirm = $('#r-confirm').value;
      if (pw !== confirm) { errEl.textContent = 'Kata sandi tidak cocok'; errEl.classList.remove('hidden'); return; }
      btn.disabled = true; btn.innerHTML = 'Memproses...';
      try {
        await API.register($('#r-name').value.trim(), $('#r-email').value.trim(), pw);
        toast('Registrasi berhasil!'); location.hash = '#/onboarding';
      } catch (err) { errEl.textContent = err.message; errEl.classList.remove('hidden'); btn.disabled = false; btn.innerHTML = `Daftar Sekarang ${Icons.arrowRight(20, '#fff')}`; }
    };
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: ONBOARDING
  // ════════════════════════════════════════════════════════════════
  function renderOnboarding() {
    app().innerHTML = `
    <div class="relative flex h-auto min-h-screen w-full flex-col overflow-x-hidden">
      <header class="flex items-center justify-between px-6 py-6 lg:px-40">
        <div class="flex items-center gap-2 text-primary">
          ${Icons.layers(28, '#2ec2b3')}
          <h2 class="text-slate-100 text-xl font-bold leading-tight tracking-tight">HabitFlow</h2>
        </div>
        <a href="#/dashboard" class="text-slate-400 hover:text-primary transition-colors text-sm font-semibold">Lewati</a>
      </header>
      <main class="flex flex-1 flex-col items-center justify-center px-6 lg:px-40 pb-12">
        <div class="flex flex-col max-w-[480px] w-full gap-8">
          <div class="flex flex-col gap-8 text-center">
            <div class="relative aspect-square w-full max-w-[320px] mx-auto bg-primary/10 rounded-full flex items-center justify-center overflow-hidden">
              <div class="w-full h-full p-8 flex flex-col justify-center items-center">
                <div class="w-full h-32 flex items-end gap-2 justify-center">
                  <div class="w-6 bg-primary/20 h-12 rounded-t-lg"></div>
                  <div class="w-6 bg-primary/40 h-20 rounded-t-lg"></div>
                  <div class="w-6 bg-primary/60 h-16 rounded-t-lg"></div>
                  <div class="w-6 bg-primary/80 h-24 rounded-t-lg"></div>
                  <div class="w-6 bg-primary h-28 rounded-t-lg"></div>
                </div>
                <div class="mt-4 flex items-center gap-2 text-primary font-bold">
                  ${Icons.flame(24, '#2ec2b3')}
                  <span class="text-2xl">15 Hari Beruntun</span>
                </div>
              </div>
            </div>
            <div class="flex flex-col gap-3">
              <h1 class="text-2xl lg:text-3xl font-bold text-slate-100">Lihat Konsistensimu</h1>
              <p class="text-slate-400 text-base lg:text-lg leading-relaxed">Pantau kemajuan dan pertahankan ritme kebiasaanmu dengan visualisasi yang memotivasi.</p>
            </div>
          </div>
          <div class="flex items-center justify-center gap-3">
            <div class="size-2.5 rounded-full bg-slate-700"></div>
            <div class="size-2.5 rounded-full bg-slate-700"></div>
            <div class="w-8 h-2.5 rounded-full bg-primary"></div>
          </div>
          <div class="flex flex-col gap-4 pt-4">
            <a href="#/dashboard" class="flex w-full cursor-pointer items-center justify-center overflow-hidden rounded-xl h-14 px-5 bg-primary text-white text-lg font-bold hover:opacity-90 transition-opacity shadow-lg shadow-primary/20">Mulai Sekarang</a>
          </div>
        </div>
      </main>
      <section class="hidden lg:flex px-40 py-12 border-t border-slate-800 gap-8 justify-center">
        <div class="flex flex-1 flex-col gap-4 max-w-[280px] opacity-40">
          <div class="aspect-video bg-primary/10 rounded-lg flex items-center justify-center">${Icons.listChecks(40, '#2ec2b3')}</div>
          <div><p class="text-slate-100 text-sm font-bold">Tambahkan Kebiasaan</p><p class="text-slate-400 text-xs">Mulai dari satu kebiasaan kecil.</p></div>
        </div>
        <div class="flex flex-1 flex-col gap-4 max-w-[280px] opacity-40">
          <div class="aspect-video bg-primary/10 rounded-lg flex items-center justify-center">${Icons.tapHand(40, '#2ec2b3')}</div>
          <div><p class="text-slate-100 text-sm font-bold">Satu Tap untuk Check-off</p><p class="text-slate-400 text-xs">Selesaikan tugas harianmu dengan mudah.</p></div>
        </div>
        <div class="flex flex-1 flex-col gap-4 max-w-[280px]">
          <div class="aspect-video bg-primary/20 rounded-lg flex items-center justify-center border-2 border-primary">${Icons.chartLine(40, '#2ec2b3')}</div>
          <div><p class="text-slate-100 text-sm font-bold">Lihat Konsistensimu</p><p class="text-slate-400 text-xs">Pantau kemajuan setiap harinya.</p></div>
        </div>
      </section>
    </div>`;
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: DASHBOARD (HARI INI)
  // ════════════════════════════════════════════════════════════════
  async function renderDashboard() {
    const user = API.getUser();
    const headerHTML = `
      <div>
        <h1 class="text-xl font-bold">Hai, ${escapeHtml(user?.name || 'User')} <span class="inline-block align-middle">${Icons.handWave(22, '#fbbf24')}</span></h1>
        <p class="text-sm text-slate-500">${formatDateID()}</p>
      </div>
      <div class="flex items-center gap-4">
        <button class="size-10 rounded-full flex items-center justify-center bg-slate-800 text-slate-400 tap-target">${Icons.bell(20)}</button>
      </div>`;
    const contentHTML = `
      <div id="progress-summary"><div class="skeleton h-28 w-full"></div></div>
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-bold">Daftar Habit</h2>
        <a href="#/habits" class="text-primary text-sm font-semibold">Lihat Semua</a>
      </div>
      <div id="habit-list" class="space-y-4 pb-20 md:pb-6">${skeleton(3)}</div>`;
    app().innerHTML = appShell('dashboard', headerHTML, contentHTML);
    try {
      const res = await API.getToday();
      const habits = res.data.habits || [];
      const doneCount = habits.filter(h => h.is_done_today).length;
      const total = habits.length;
      const pct = total ? Math.round(doneCount / total * 100) : 0;
      const remaining = total - doneCount;
      const allDone = total > 0 && doneCount === total;
      $('#progress-summary').innerHTML = total > 0
        ? `<div class="bg-primary rounded-xl p-6 text-slate-950 shadow-lg shadow-primary/10 flex items-center justify-between">
             <div class="space-y-2">
               <h3 class="text-lg font-medium opacity-90">Progress Hari Ini</h3>
               <p class="text-3xl font-bold"><span class="progress-done">${doneCount}</span>/<span class="progress-total">${total}</span> <span class="text-lg font-normal opacity-75">Selesai</span></p>
               <p class="progress-msg text-sm bg-black/20 inline-block px-3 py-1 rounded-full">${allDone ? `Semua habit selesai! ${Icons.partyPopper(16, '#fff')}` : `${remaining} habit lagi untuk mencapai target!`}</p>
             </div>
             <div class="relative size-24 flex items-center justify-center">
               <svg class="size-full -rotate-90" viewBox="0 0 36 36">
                 <path class="text-white/20" d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" fill="none" stroke="currentColor" stroke-width="3"></path>
                 <path class="progress-ring-fill text-white" d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831" fill="none" stroke="currentColor" stroke-dasharray="${pct}, 100" stroke-linecap="round" stroke-width="3"></path>
               </svg>
               <span class="progress-pct absolute text-xl font-bold">${pct}%</span>
             </div>
           </div>`
        : `<div class="bg-slate-900/50 border border-slate-800 rounded-xl p-6 text-center"><p class="text-slate-400 text-sm">Belum ada habit. Tambahkan yang pertama! ${Icons.sprout(16, '#2ec2b3')}</p></div>`;
      const list = $('#habit-list');
      if (habits.length === 0) { list.innerHTML = `<div class="text-center text-slate-500 py-16">${Icons.emptyHabits()}<p class="text-sm mt-4">Mulai dengan menambahkan habit pertamamu</p></div>`; return; }
      list.innerHTML = habits.map(h => habitCard(h)).join('');
      bindHabitCards();
    } catch (err) { $('#habit-list').innerHTML = `<div class="text-center text-red-400 py-8 text-sm">${escapeHtml(err.message)}</div>`; }
  }

  function habitCard(h) {
    const done = h.is_done_today;
    const colors = categoryColor(h.category);
    const streakClass = h.current_streak > 0 ? 'text-orange-500 font-medium' : 'text-slate-500';
    const fireHTML = h.current_streak > 0 ? `<span class="fire-pulse inline-flex">${Icons.flame(16, '#f97316')}</span>` : Icons.flame(16, '#64748b');
    const badgeHTML = Icons.streakBadge(h.current_streak, 16);
    return `
    <div class="habit-card ${done ? 'done bg-slate-900/20 border border-transparent' : 'bg-slate-900/50 border border-slate-800 hover:border-primary/30'} p-4 rounded-xl shadow-sm flex items-center justify-between group transition-colors tap-target" data-id="${h.habit_id}" data-streak="${h.current_streak}">
      <div class="flex items-center gap-4">
        <div class="size-12 rounded-xl ${colors.bg} ${colors.text} flex items-center justify-center">${Icons.category(h.category, 24)}</div>
        <div>
          <h4 class="habit-name font-bold text-slate-100 ${done ? 'line-through decoration-slate-400' : ''}">${escapeHtml(h.name)}</h4>
          <p class="text-sm text-slate-500 flex items-center gap-1">
            <span class="streak-info flex items-center gap-1 ${streakClass}">
              ${fireHTML} <span class="streak-count">${h.current_streak}</span> hari ${badgeHTML}
            </span> • ${categoryLabel(h.category)}
          </p>
        </div>
      </div>
      ${done
        ? `<button type="button" class="check-btn size-12 rounded-full bg-primary text-white flex items-center justify-center shrink-0 transition-all" data-habit-id="${h.habit_id}" data-done="true" role="button" tabindex="0">${Icons.checkCheck(24, '#fff')}</button>`
        : `<button type="button" class="check-btn size-12 rounded-full border-2 border-primary text-primary flex items-center justify-center hover:bg-primary/5 transition-all shrink-0" data-habit-id="${h.habit_id}" data-done="false" tabindex="0">${Icons.check(24)}</button>`
      }
    </div>`;
  }

  function updateProgressSummary(change) {
    const el = $('#progress-summary');
    if (!el) return;
    const countEl = el.querySelector('.progress-done');
    const totalEl = el.querySelector('.progress-total');
    const pctEl = el.querySelector('.progress-pct');
    const msgEl = el.querySelector('.progress-msg');
    const ringEl = el.querySelector('.progress-ring-fill');
    if (!countEl || !totalEl) return;
    const done = parseInt(countEl.textContent, 10) + change;
    const total = parseInt(totalEl.textContent, 10);
    const pct = total ? Math.round(done / total * 100) : 0;
    const remaining = total - done;
    const allDone = total > 0 && done >= total;
    countEl.textContent = done;
    if (pctEl) pctEl.textContent = pct + '%';
    if (msgEl) msgEl.innerHTML = allDone ? `Semua habit selesai! ${Icons.partyPopper(16, '#fff')}` : `${remaining} habit lagi untuk mencapai target!`;
    if (ringEl) ringEl.setAttribute('stroke-dasharray', `${pct}, 100`);
  }

  // ── SVG icons for daily summary ──
  const trendUpSVG = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m18 15-6-6-6 6"/></svg>';
  const trendDownSVG = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>';
  const trendSameSVG = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/></svg>';
  const flameSVG = '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8.5 14.5A2.5 2.5 0 0 0 11 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 1 1-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 0 0 2.5 2.5z"/></svg>';
  const compareSVG = '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 3h5v5"/><path d="M8 3H3v5"/><path d="M12 22V8"/><path d="m21 3-9 9"/><path d="M3 3l9 9"/></svg>';

  function bindHabitCards() {
    $$('.check-btn').forEach(btn => {
      btn.onclick = async () => {
        if (btn.dataset.busy === '1') return;
        btn.dataset.busy = '1';
        const id = btn.dataset.habitId;
        const wasDone = btn.dataset.done === 'true';
        const card = btn.closest('.habit-card');

        if (wasDone) {
          // ── UNDO ──
          if (!confirm('Batalkan checkin hari ini?')) { btn.dataset.busy = '0'; return; }
          const prevStreak = parseInt(card?.dataset.streak || '0', 10);
          // Optimistic UI
          btn.dataset.done = 'false';
          btn.className = 'check-btn size-12 rounded-full border-2 border-primary text-primary flex items-center justify-center hover:bg-primary/5 transition-all shrink-0';
          btn.innerHTML = Icons.check(24);
          if (card) {
            card.classList.remove('done', 'bg-slate-900/20', 'border-transparent');
            card.classList.add('bg-slate-900/50', 'border-slate-800', 'hover:border-primary/30');
            const nm = card.querySelector('.habit-name');
            if (nm) nm.classList.remove('line-through', 'decoration-slate-400');
          }
          updateProgressSummary(-1);
          try {
            await API.undoCheckin(id);
            toast('Checkin dibatalkan', 'info');
            // Update streak from server — re-fetch today data silently
            try {
              const res = await API.getToday();
              const h = (res.data.habits || []).find(x => String(x.habit_id) === String(id));
              if (h && card) {
                card.dataset.streak = h.current_streak;
                const sc = card.querySelector('.streak-count');
                if (sc) sc.textContent = h.current_streak;
                const si = card.querySelector('.streak-info');
                if (si) {
                  si.className = 'streak-info ' + (h.current_streak > 0 ? 'text-orange-500 font-medium' : 'text-slate-500');
                }
              }
            } catch (_) {}
          } catch (err) {
            // Rollback
            toast(err.message, 'error');
            btn.dataset.done = 'true';
            btn.className = 'check-btn size-12 rounded-full bg-primary text-white flex items-center justify-center shrink-0 transition-all';
            btn.innerHTML = Icons.checkCheck(24, '#fff');
            if (card) {
              card.classList.add('done', 'bg-slate-900/20', 'border-transparent');
              card.classList.remove('bg-slate-900/50', 'border-slate-800', 'hover:border-primary/30');
              const nm = card.querySelector('.habit-name');
              if (nm) nm.classList.add('line-through', 'decoration-slate-400');
              card.dataset.streak = prevStreak;
              const sc = card.querySelector('.streak-count');
              if (sc) sc.textContent = prevStreak;
            }
            updateProgressSummary(1);
          }
          btn.dataset.busy = '0';
        } else {
          // ── CHECK-OFF ──
          const prevStreak = parseInt(card?.dataset.streak || '0', 10);
          // Optimistic UI
          btn.dataset.done = 'true';
          btn.className = 'check-btn size-12 rounded-full bg-primary text-white flex items-center justify-center shrink-0 check-pop transition-all';
          btn.innerHTML = Icons.checkCheck(24, '#fff');
          if (card) {
            card.classList.add('done', 'bg-slate-900/20', 'border-transparent');
            card.classList.remove('bg-slate-900/50', 'border-slate-800', 'hover:border-primary/30');
            const nm = card.querySelector('.habit-name');
            if (nm) nm.classList.add('line-through', 'decoration-slate-400');
            // Optimistic streak +1
            const newStreak = prevStreak + 1;
            card.dataset.streak = newStreak;
            const sc = card.querySelector('.streak-count');
            if (sc) sc.textContent = newStreak;
            const si = card.querySelector('.streak-info');
            if (si) si.className = 'streak-info text-orange-500 font-medium';
          }
          updateProgressSummary(1);
          try {
            const res = await API.checkin(id);
            toast('Habit berhasil dicheckin!');
            // Correct streak from server
            if (res.data && card) {
              card.dataset.streak = res.data.current_streak;
              const sc = card.querySelector('.streak-count');
              if (sc) sc.textContent = res.data.current_streak;
            }
          } catch (err) {
            // Rollback
            toast(err.message, 'error');
            btn.dataset.done = 'false';
            btn.className = 'check-btn size-12 rounded-full border-2 border-primary text-primary flex items-center justify-center hover:bg-primary/5 transition-all shrink-0';
            btn.innerHTML = Icons.check(24);
            if (card) {
              card.classList.remove('done', 'bg-slate-900/20', 'border-transparent');
              card.classList.add('bg-slate-900/50', 'border-slate-800', 'hover:border-primary/30');
              const nm = card.querySelector('.habit-name');
              if (nm) nm.classList.remove('line-through', 'decoration-slate-400');
              card.dataset.streak = prevStreak;
              const sc = card.querySelector('.streak-count');
              if (sc) sc.textContent = prevStreak;
              const si = card.querySelector('.streak-info');
              if (si) si.className = 'streak-info ' + (prevStreak > 0 ? 'text-orange-500 font-medium' : 'text-slate-500');
            }
            updateProgressSummary(-1);
          }
          btn.dataset.busy = '0';
        }
      };
    });
    $$('.habit-card').forEach(card => { card.addEventListener('dblclick', () => { location.hash = '#/habits/' + card.dataset.id + '/edit'; }); });
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: HABIT LIST (KEBIASAAN SAYA)
  // ════════════════════════════════════════════════════════════════
  async function renderHabits() {
    const headerHTML = `
      <div>
        <h2 class="text-xl font-bold">Kebiasaan Saya</h2>
        <p class="text-sm text-slate-400">Pantau dan kelola semua target harian Anda.</p>
      </div>
      <div class="flex items-center gap-4">
        <div class="relative hidden md:block">
          <span class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400">${Icons.search(20)}</span>
          <input id="habit-search" class="pl-10 pr-4 py-2 bg-slate-800 border-none rounded-lg text-sm w-64 focus:ring-2 focus:ring-primary placeholder:text-slate-500" placeholder="Cari kebiasaan..." type="text">
        </div>
        <div class="size-10 rounded-full bg-slate-700 overflow-hidden border-2 border-slate-800 flex items-center justify-center text-primary font-bold text-sm">
          ${(API.getUser()?.name || 'U')[0].toUpperCase()}
        </div>
      </div>`;
    const contentHTML = `
      <div class="max-w-6xl mx-auto">
        <div class="flex items-center justify-between mb-8">
          <div class="flex p-1 bg-slate-800/50 rounded-xl w-64">
            <button id="filter-active" class="flex-1 py-2 px-4 rounded-lg text-sm font-medium bg-slate-700 shadow-sm text-primary">Aktif</button>
            <button id="filter-all" class="flex-1 py-2 px-4 rounded-lg text-sm font-medium text-slate-400 hover:text-slate-200">Semua</button>
          </div>
          <a href="#/habits/new" class="md:hidden flex items-center gap-2 px-4 py-2 bg-primary text-slate-900 rounded-xl text-sm font-semibold shadow-lg shadow-primary/20">
            ${Icons.plus(18)} Tambah
          </a>
        </div>
        <div id="habits-table">${skeleton(5)}</div>
      </div>
      <a href="#/habits/new" class="hidden md:flex fixed bottom-8 right-8 items-center gap-2 px-6 py-3 bg-primary text-slate-900 rounded-xl font-semibold shadow-lg shadow-primary/20 hover:opacity-90 transition-all z-20 tap-target">
        ${Icons.plus(20)} Tambah Kebiasaan
      </a>`;
    app().innerHTML = appShell('habits', headerHTML, contentHTML, { fab: false });
    let showAll = false;
    async function loadHabits() {
      try {
        const res = await API.getHabits();
        let habits = res.data || [];
        const totalAll = habits.length;
        if (!showAll) habits = habits.filter(h => h.is_active);
        const container = $('#habits-table');
        if (habits.length === 0) { container.innerHTML = `<div class="text-center py-16 text-slate-500">${Icons.emptyHabits()}<p class="mt-4">Belum ada kebiasaan. Tambahkan yang pertama!</p></div>`; return; }
        container.innerHTML = `
          <div class="bg-slate-800/40 rounded-2xl border border-slate-800 overflow-hidden">
            <table class="w-full text-left border-collapse">
              <thead><tr class="bg-slate-800/50 border-b border-slate-800">
                <th class="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Ikon</th>
                <th class="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500">Nama Kebiasaan</th>
                <th class="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500 hidden md:table-cell">Streak Terlama</th>
                <th class="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500 hidden sm:table-cell">Status</th>
                <th class="px-6 py-4 text-xs font-semibold uppercase tracking-wider text-slate-500 text-right">Aksi</th>
              </tr></thead>
              <tbody class="divide-y divide-slate-800">${habits.map(h => habitRow(h)).join('')}</tbody>
            </table>
          </div>
          <div class="mt-6 flex items-center justify-between text-sm text-slate-400 px-2">
            <p>Menampilkan ${habits.length} dari ${totalAll} kebiasaan</p>
          </div>`;
      } catch (err) { $('#habits-table').innerHTML = `<div class="text-center text-red-400 py-8">${escapeHtml(err.message)}</div>`; }
    }
    function habitRow(h) {
      const colors = categoryColor(h.category);
      return `
      <tr class="habit-row hover:bg-slate-800/40 transition-colors">
        <td class="px-6 py-5"><div class="w-10 h-10 rounded-xl ${colors.bg} ${colors.text} flex items-center justify-center">${Icons.category(h.category, 20)}</div></td>
        <td class="px-6 py-5">
          <p class="font-semibold text-slate-100">${escapeHtml(h.name)}</p>
          <p class="text-xs text-slate-500">Harian${h.notify_time ? ' • ' + h.notify_time : ''}</p>
        </td>
        <td class="px-6 py-5 hidden md:table-cell">
          <div class="flex items-center gap-2">${Icons.flame(18, '#f97316')}<span class="text-sm font-medium">${h.longest_streak || 0} Hari</span></div>
        </td>
        <td class="px-6 py-5 hidden sm:table-cell">
          <span class="px-2.5 py-1 rounded-full text-[11px] font-bold uppercase ${h.is_active ? 'bg-emerald-900/30 text-emerald-400' : 'bg-slate-800 text-slate-400'}">${h.is_active ? 'Aktif' : 'Tidak Aktif'}</span>
        </td>
        <td class="px-6 py-5 text-right">
          <a href="#/habits/${h.id}/edit" class="text-slate-400 hover:text-primary transition-colors">${Icons.pencil(18)}</a>
        </td>
      </tr>`;
    }
    await loadHabits();
    const filterActive = $('#filter-active'), filterAll = $('#filter-all');
    if (filterActive) filterActive.onclick = () => { showAll = false; filterActive.className = 'flex-1 py-2 px-4 rounded-lg text-sm font-medium bg-slate-700 shadow-sm text-primary'; filterAll.className = 'flex-1 py-2 px-4 rounded-lg text-sm font-medium text-slate-400 hover:text-slate-200'; loadHabits(); };
    if (filterAll) filterAll.onclick = () => { showAll = true; filterAll.className = 'flex-1 py-2 px-4 rounded-lg text-sm font-medium bg-slate-700 shadow-sm text-primary'; filterActive.className = 'flex-1 py-2 px-4 rounded-lg text-sm font-medium text-slate-400 hover:text-slate-200'; loadHabits(); };
    const search = $('#habit-search');
    if (search) search.oninput = () => {
      const q = search.value.toLowerCase();
      $$('.habit-row').forEach(row => { row.style.display = row.textContent.toLowerCase().includes(q) ? '' : 'none'; });
    };
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: HABIT FORM (ADD / EDIT)
  // ════════════════════════════════════════════════════════════════
  async function renderHabitForm(editId) {
    let habit = null;
    if (editId) { try { const res = await API.getHabit(editId); habit = res.data; } catch { toast('Habit tidak ditemukan', 'error'); location.hash = '#/habits'; return; } }
    const categories = ['health', 'fitness', 'learning', 'productivity', 'mindfulness', 'finance', 'social', 'general'];
    const initialNotify = habit?.notify_time || '';
    const initialIsPreset = isPresetNotifyTime(initialNotify);
    const initialUseCustomTime = !!initialNotify && !initialIsPreset;
    const initialCustomTimeValue = initialUseCustomTime ? initialNotify : '';
    const headerHTML = `
      <div class="flex items-center gap-3">
        <a href="#/habits" class="size-10 flex items-center justify-center rounded-xl bg-slate-800 text-slate-400 tap-target">${Icons.arrowLeft(20)}</a>
        <h1 class="text-lg font-bold text-white">${habit ? 'Edit Habit' : 'Habit Baru'}</h1>
      </div>
      <div></div>`;
    const contentHTML = `
      <form id="habit-form" class="max-w-lg space-y-5">
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-slate-300">Nama Habit</label>
          <div class="group relative flex w-full items-center">
            <div class="absolute left-3 text-slate-400 group-focus-within:text-primary transition-colors">${Icons.pencil(20)}</div>
            <input id="f-name" type="text" value="${habit ? escapeHtml(habit.name) : ''}" placeholder="contoh: Membaca 10 halaman" required maxlength="100"
              class="flex h-12 w-full rounded-xl border border-slate-700 bg-slate-800/50 px-10 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all">
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-slate-300">Kategori</label>
          <div class="relative w-full" id="category-dropdown">
            <input type="hidden" id="f-category" value="${habit?.category || 'general'}">
            <button type="button" id="category-trigger"
              class="w-full flex items-center justify-between px-4 py-3 h-12
                     bg-slate-800/50 border border-slate-700 rounded-xl
                     text-sm text-slate-100 hover:border-primary/50
                     focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary
                     transition-all duration-200"
              aria-haspopup="listbox" aria-expanded="false">
              <span class="flex items-center gap-3" id="category-display">
                ${categorySVG(habit?.category || 'general')}
                <span>${categoryLabel(habit?.category || 'general')}</span>
              </span>
              <span id="category-chevron" class="transition-transform duration-200 text-slate-400">${chevronDownSVG}</span>
            </button>
            <div id="category-panel"
              class="absolute top-full left-0 right-0 mt-2 py-2
                     bg-slate-900/95 backdrop-blur-md border border-slate-700 rounded-xl shadow-2xl
                     z-50 hidden dropdown-panel
                     max-h-[60vh] overflow-y-auto hide-scrollbar"
              role="listbox" tabindex="-1">
              ${categories.map(c => `
                <button type="button" class="category-option w-full flex items-center gap-3 px-4 py-3
                                              text-sm text-slate-200 hover:bg-primary/10
                                              focus:bg-primary/10 focus:outline-none
                                              transition-colors"
                        data-value="${c}" role="option">
                  ${categorySVG(c)}
                  <span>${categoryLabel(c)}</span>
                </button>
              `).join('')}
            </div>
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-slate-300">Waktu Notifikasi <span class="text-slate-500 font-normal">(opsional)</span></label>
          <div class="relative w-full" id="time-dropdown">
            <input type="hidden" id="f-notify" value="${initialUseCustomTime ? '' : initialNotify}">
            <button type="button" id="time-trigger"
              class="w-full flex items-center justify-between px-4 py-3 h-12
                     bg-slate-800/50 border border-slate-700 rounded-xl
                     text-sm text-slate-100 hover:border-primary/50
                     focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary
                     transition-all duration-200"
              aria-haspopup="listbox" aria-expanded="false">
              <span class="flex items-center gap-3" id="time-display">
                ${notifyTimeDisplayHTML(habit?.notify_time || '')}
              </span>
              <span id="time-chevron" class="transition-transform duration-200 text-slate-400">${chevronDownSVG}</span>
            </button>
            <div id="time-panel"
              class="absolute top-full left-0 right-0 mt-2 py-2
                     bg-slate-900/95 backdrop-blur-md border border-slate-700 rounded-xl shadow-2xl
                     z-50 hidden dropdown-panel
                     max-h-[60vh] overflow-y-auto hide-scrollbar"
              role="listbox" tabindex="-1">
              ${buildTimeOptions(initialUseCustomTime ? '' : initialNotify)}
            </div>
          </div>
          <label class="inline-flex items-center gap-2 mt-1 text-xs text-slate-400">
            <input type="checkbox" id="f-custom-time-enabled" class="rounded border-slate-600 bg-slate-800/60 text-primary focus:ring-primary/40" ${initialUseCustomTime ? 'checked' : ''}>
            Atur waktu sendiri
          </label>
          <div id="custom-time-wrap" class="${initialUseCustomTime ? '' : 'hidden'}">
            <div class="relative w-full" id="custom-time-dropdown">
              <input id="f-custom-time" type="hidden" value="${initialCustomTimeValue}">
              <button type="button" id="custom-time-trigger"
                class="w-full flex items-center justify-between px-4 py-3 h-12
                       bg-slate-800/50 border border-slate-700 rounded-xl
                       text-sm text-slate-100 hover:border-primary/50
                       focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary
                       transition-all duration-200"
                aria-haspopup="listbox" aria-expanded="false"
                aria-label="Pilih waktu sendiri">
                <span class="flex items-center gap-3" id="custom-time-display">${customTimeDisplayHTML(initialCustomTimeValue)}</span>
                <span id="custom-time-chevron" class="text-slate-400 transition-transform duration-200">${chevronDownSVG}</span>
              </button>
              <div id="custom-time-panel"
                class="absolute top-full left-0 right-0 mt-2 py-2
                       bg-slate-900/95 backdrop-blur-md border border-slate-700 rounded-xl shadow-2xl
                       z-50 hidden dropdown-panel
                       max-h-[60vh] overflow-y-auto hide-scrollbar"
                role="listbox" tabindex="-1">
                ${buildCustomTimeOptions(initialCustomTimeValue)}
              </div>
            </div>
            <p class="text-[11px] text-slate-500 mt-2">Format 24 jam, contoh 21:30</p>
          </div>
        </div>
        <button type="submit" id="btn-save" class="flex w-full items-center justify-center rounded-xl bg-primary h-12 px-4 text-sm font-bold text-white transition-all hover:bg-primary/90 active:scale-[0.98] shadow-lg shadow-primary/20">${habit ? 'Simpan Perubahan' : 'Buat Habit'}</button>
        <a href="#/habits" class="block text-center py-3 text-slate-500 text-sm font-medium hover:text-slate-300 transition tap-target">Batal</a>
        ${habit ? `<div class="pt-2 border-t border-slate-800"><button type="button" id="btn-delete" class="w-full py-3 bg-slate-900/50 border border-red-500/30 text-red-400 font-semibold rounded-xl hover:bg-red-900/20 transition tap-target">Hapus Habit</button></div>` : ''}
      </form>`;
    app().innerHTML = appShell('habits', headerHTML, contentHTML, { fab: false });
    initCategoryDropdown();
    initTimeDropdown();

    const hiddenNotifyInput = $('#f-notify');
    const customToggle = $('#f-custom-time-enabled');
    const customWrap = $('#custom-time-wrap');
    const customInput = $('#f-custom-time');
    const customTrigger = $('#custom-time-trigger');
    const customPanel = $('#custom-time-panel');
    const customChevron = $('#custom-time-chevron');
    const customDisplay = $('#custom-time-display');
    const timeDisplay = $('#time-display');

    function updateCustomTimeDisplay() {
      if (customDisplay) customDisplay.innerHTML = customTimeDisplayHTML(customInput?.value || '');
    }

    function syncCustomOptionSelection() {
      if (!customPanel || !customInput) return;
      const current = customInput.value || '';
      const options = $$('.custom-time-option', customPanel);
      options.forEach(o => {
        const selected = o.dataset.value === current;
        o.classList.toggle('bg-primary/10', selected);
        o.classList.toggle('text-primary', selected);
        o.classList.toggle('border-primary/30', selected);
        o.classList.toggle('custom-time-option-selected', selected);
        const check = o.querySelector('.time-option-check');
        if (check) {
          check.classList.toggle('opacity-100', selected);
          check.classList.toggle('opacity-0', !selected);
        }
      });
    }

    function initCustomTimeDropdown() {
      if (!customTrigger || !customPanel || !customInput) return;
      let isOpen = false;
      const options = $$('.custom-time-option', customPanel);
      const hourInput = $('#custom-hour-input', customPanel);
      const minuteInput = $('#custom-minute-input', customPanel);
      const applyButton = $('#custom-time-apply', customPanel);
      let focusIdx = -1;

      function normalizeTime(v) {
        const m = /^([01]?\d|2[0-3]):([0-5]\d)$/.exec((v || '').trim());
        if (!m) return '';
        return `${m[1].padStart(2, '0')}:${m[2]}`;
      }

      function setManualInputsFromValue(v) {
        const normalized = normalizeTime(v);
        if (!normalized) {
          if (hourInput) hourInput.value = '';
          if (minuteInput) minuteInput.value = '';
          return;
        }
        const [hh, mm] = normalized.split(':');
        if (hourInput) hourInput.value = hh;
        if (minuteInput) minuteInput.value = mm;
      }

      function applyManualValue() {
        const hh = (hourInput?.value || '').replace(/\D/g, '');
        const mm = (minuteInput?.value || '').replace(/\D/g, '');
        if (hh === '' || mm === '') {
          toast('Lengkapi jam dan menit dulu', 'error');
          return;
        }
        const h = Number(hh);
        const m = Number(mm);
        if (Number.isNaN(h) || h < 0 || h > 23 || Number.isNaN(m) || m < 0 || m > 59) {
          toast('Format waktu tidak valid (00:00 - 23:59)', 'error');
          return;
        }
        customInput.value = `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`;
        syncCustomOptionSelection();
        applyCustomTimeState();
        close();
        customTrigger.focus();
      }

      function open() {
        isOpen = true;
        customPanel.classList.remove('hidden');
        customPanel.classList.add('dropdown-enter');
        customChevron?.classList.add('rotate-180');
        customTrigger.setAttribute('aria-expanded', 'true');
        focusIdx = -1;
        setManualInputsFromValue(customInput.value || '');
        setTimeout(() => hourInput?.focus(), 50);
      }

      function close() {
        isOpen = false;
        customPanel.classList.add('dropdown-exit');
        customChevron?.classList.remove('rotate-180');
        customTrigger.setAttribute('aria-expanded', 'false');
        setTimeout(() => {
          customPanel.classList.add('hidden');
          customPanel.classList.remove('dropdown-enter', 'dropdown-exit');
        }, 150);
        focusIdx = -1;
      }

      function selectOption(opt) {
        customInput.value = opt.dataset.value || '';
        syncCustomOptionSelection();
        applyCustomTimeState();
        close();
        customTrigger.focus();
      }

      function focusOption(idx) {
        if (options.length === 0) return;
        if (idx < 0) idx = options.length - 1;
        if (idx >= options.length) idx = 0;
        focusIdx = idx;
        options.forEach(o => o.classList.remove('time-option-focus'));
        options[focusIdx].classList.add('time-option-focus');
        options[focusIdx].scrollIntoView({ block: 'nearest' });
      }

      customTrigger.addEventListener('click', (e) => { e.preventDefault(); isOpen ? close() : open(); });
      options.forEach(opt => opt.addEventListener('click', (e) => { e.preventDefault(); selectOption(opt); }));
      if (applyButton) applyButton.addEventListener('click', (e) => { e.preventDefault(); applyManualValue(); });
      document.addEventListener('click', (e) => {
        if (isOpen && !customTrigger.contains(e.target) && !customPanel.contains(e.target)) close();
      });

      const clampNumericInput = (el) => {
        if (!el) return;
        el.addEventListener('input', () => {
          el.value = el.value.replace(/\D/g, '').slice(0, 2);
        });
      };
      clampNumericInput(hourInput);
      clampNumericInput(minuteInput);

      [hourInput, minuteInput].forEach(el => {
        if (!el) return;
        el.addEventListener('keydown', (e) => {
          if (e.key === 'Enter') {
            e.preventDefault();
            applyManualValue();
          }
        });
      });

      customTrigger.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
          e.preventDefault();
          if (!isOpen) open();
          if (options.length > 0) focusOption(e.key === 'ArrowDown' ? 0 : options.length - 1);
        }
        if (e.key === 'Escape' && isOpen) { e.preventDefault(); close(); }
        if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); isOpen ? close() : open(); }
      });

      customPanel.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowDown') { e.preventDefault(); if (options.length > 0) focusOption(focusIdx + 1); }
        else if (e.key === 'ArrowUp') { e.preventDefault(); if (options.length > 0) focusOption(focusIdx - 1); }
        else if ((e.key === 'Enter' || e.key === ' ') && focusIdx >= 0 && options.length > 0) { e.preventDefault(); selectOption(options[focusIdx]); }
        else if (e.key === 'Escape') { e.preventDefault(); close(); customTrigger.focus(); }
        else if (e.key === 'Tab') { close(); }
      });
    }

    updateCustomTimeDisplay();
    syncCustomOptionSelection();
    initCustomTimeDropdown();

    function applyCustomTimeState() {
      const useCustom = !!customToggle?.checked;
      if (customWrap) customWrap.classList.toggle('hidden', !useCustom);
      if (!useCustom) {
        if (customInput) customInput.value = '';
        if (hiddenNotifyInput) hiddenNotifyInput.value = '';
        if (timeDisplay) timeDisplay.innerHTML = notifyTimeDisplayHTML('');
        return;
      }

      const value = customInput?.value || '';
      if (hiddenNotifyInput) hiddenNotifyInput.value = value;
      if (timeDisplay) timeDisplay.innerHTML = notifyTimeDisplayHTML(value);
      updateCustomTimeDisplay();
      syncCustomOptionSelection();
    }

    if (customToggle && customInput) {
      customToggle.addEventListener('change', () => {
        applyCustomTimeState();
        if (customToggle.checked) {
          customTrigger?.focus();
          if (!customInput.value) {
            customInput.value = '21:30';
            applyCustomTimeState();
          }
        }
      });

      customInput.addEventListener('change', applyCustomTimeState);
      customInput.addEventListener('input', applyCustomTimeState);
    }

    if (hiddenNotifyInput && customToggle) {
      hiddenNotifyInput.addEventListener('change', () => {
        if (hiddenNotifyInput.value !== '' && isPresetNotifyTime(hiddenNotifyInput.value)) {
          customToggle.checked = false;
          if (customWrap) customWrap.classList.add('hidden');
          if (customInput) customInput.value = '';
        }
      });
    }

    $('#habit-form').onsubmit = async (e) => {
      e.preventDefault(); const btn = $('#btn-save'); btn.disabled = true; btn.textContent = 'Menyimpan...';

      if (customToggle?.checked) {
        const customVal = customInput?.value || '';
        if (!customVal) {
          toast('Isi waktu notifikasi custom terlebih dulu', 'error');
          btn.disabled = false;
          btn.textContent = habit ? 'Simpan Perubahan' : 'Buat Habit';
          return;
        }
        hiddenNotifyInput.value = customVal;
      }

      // Send empty string instead of null so backend update path treats this as an explicit notify_time change.
      const body = { name: $('#f-name').value.trim(), category: $('#f-category').value, notify_time: $('#f-notify').value || '' };
      try {
        if (habit) { await API.updateHabit(editId, body); toast('Habit berhasil diupdate'); }
        else { await API.createHabit(body); toast('Habit baru berhasil dibuat!'); }
        location.hash = '#/habits';
      } catch (err) { toast(err.message, 'error'); btn.disabled = false; btn.textContent = habit ? 'Simpan Perubahan' : 'Buat Habit'; }
    };
    if (habit && $('#btn-delete')) {
      $('#btn-delete').onclick = async () => {
        if (!confirm(`Hapus "${escapeHtml(habit.name)}"?`)) return;
        try { await API.deleteHabit(editId); toast('Habit berhasil dihapus'); location.hash = '#/habits'; }
        catch (err) { toast(err.message, 'error'); }
      };
    }
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: WEEKLY REPORT (LAPORAN MINGGUAN)
  // ════════════════════════════════════════════════════════════════
  let reportWeekOffset = 0;

  async function renderReport() {
    const headerHTML = `
      <div class="flex items-center gap-8">
        <div class="flex items-center gap-3">
          <div class="text-primary">${Icons.fileBarChart(28, '#2ec2b3')}</div>
          <h2 class="text-lg font-bold text-white">HabitFlow</h2>
        </div>
        <div class="hidden md:flex items-center gap-9">
          <a href="#/dashboard" class="text-slate-300 text-sm font-medium hover:text-primary transition-colors">Dashboard</a>
          <span class="text-primary text-sm font-bold border-b-2 border-primary py-1">Laporan</span>
          <a href="#/habits" class="text-slate-300 text-sm font-medium hover:text-primary transition-colors">Habit</a>
          <a href="#/settings" class="text-slate-300 text-sm font-medium hover:text-primary transition-colors">Pengaturan</a>
        </div>
      </div>
      <div class="flex items-center gap-4">
        <button class="size-10 rounded-full flex items-center justify-center bg-slate-800 text-slate-400">${Icons.bell(20)}</button>
      </div>`;
    const contentHTML = `<div id="report-content" class="max-w-6xl mx-auto">${skeleton(4)}</div>`;
    app().innerHTML = appShell('report', headerHTML, contentHTML, { fab: false });
    await loadReport();
  }

  async function loadReport() {
    const content = $('#report-content');
    if (!content) return;
    content.innerHTML = skeleton(4);
    try {
      let url = '/reports/weekly';
      if (reportWeekOffset > 0) {
        const end = new Date(); end.setDate(end.getDate() - (reportWeekOffset * 7));
        const start = new Date(end); start.setDate(start.getDate() - 6);
        url += `?start=${dateStr(start)}&end=${dateStr(end)}`;
      }
      const res = await API._fetch(url);
      const r = res.data;
      const overallScore = r.score?.overall_score ?? 0;
      const habitScores = r.score?.habit_scores || [];
      const roundedScore = Math.round(overallScore);
      const scoreMsg = roundedScore >= 80 ? 'Performa Luar Biasa!' : roundedScore >= 50 ? 'Cukup Baik, Tingkatkan!' : 'Ayo Semangat Lagi!';
      const scoreDesc = roundedScore >= 80
        ? 'Kamu mempertahankan tingkat konsistensi di atas 80% minggu ini. Pertahankan ritme ini!'
        : roundedScore >= 50 ? 'Kamu sudah cukup konsisten, tapi masih bisa lebih baik lagi!'
        : 'Jangan menyerah! Mulai lagi dengan langkah kecil hari ini.';
      const ringR = 88, ringCirc = 2 * Math.PI * ringR;
      const ringOffset = ringCirc - (overallScore / 100) * ringCirc;

      let html = `
      <div class="flex flex-wrap justify-between items-end gap-3 mb-8">
        <div class="flex flex-col gap-1">
          <h1 class="text-white text-3xl md:text-4xl font-black leading-tight tracking-tight">Laporan Mingguan</h1>
          <p class="text-slate-400 text-lg">Periode: ${formatDateShort(r.start_date)} - ${formatDateShort(r.end_date)}</p>
        </div>
      </div>
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div class="lg:col-span-8 flex flex-col gap-6">
          <div class="bg-slate-900/50 rounded-xl p-8 flex flex-col md:flex-row items-center gap-8 border border-slate-800">
            <div class="relative flex items-center justify-center shrink-0">
              <svg class="size-48 transform -rotate-90">
                <circle class="text-slate-800/50" cx="96" cy="96" fill="transparent" r="${ringR}" stroke="currentColor" stroke-width="12"></circle>
                <circle class="text-primary score-ring" cx="96" cy="96" fill="transparent" r="${ringR}" stroke="currentColor" stroke-dasharray="${ringCirc.toFixed(2)}" stroke-dashoffset="${ringOffset.toFixed(2)}" stroke-width="12" stroke-linecap="round"></circle>
              </svg>
              <div class="absolute inset-0 flex flex-col items-center justify-center">
                <span class="text-5xl font-black text-white">${roundedScore}</span>
                <span class="text-primary font-bold text-sm uppercase tracking-widest">Konsisten</span>
              </div>
            </div>
            <div class="flex-1 text-center md:text-left">
              <h3 class="text-2xl font-bold text-white mb-2">${scoreMsg}</h3>
              <p class="text-slate-400 leading-relaxed mb-4">${scoreDesc}</p>
              <div class="inline-flex items-center gap-6 mt-3 text-sm text-slate-400">
                <div><span class="block text-2xl font-bold text-white">${r.total_habits || 0}</span>Habits</div>
                <div><span class="block text-2xl font-bold text-white">${r.total_checkins || 0}</span>Check-in</div>
              </div>
            </div>
          </div>`;

      if (habitScores.length > 0) {
        html += `<div class="bg-slate-900/50 rounded-xl p-6 border border-slate-800"><h3 class="text-lg font-bold text-white mb-4">Penyelesaian per Habit</h3><div class="space-y-4">`;
        habitScores.forEach(s => {
          const score = Math.round(s.score);
          html += `<div class="space-y-2"><div class="flex justify-between text-sm"><span class="font-medium">${escapeHtml(s.habit_name)}</span><span class="text-primary font-bold">${score}%</span></div><div class="h-2 bg-slate-800 rounded-full overflow-hidden"><div class="h-full bg-primary rounded-full progress-fill" style="width:${s.score}%"></div></div></div>`;
        });
        html += '</div></div>';
      }
      html += '</div><div class="lg:col-span-4 flex flex-col gap-6">';

      // Streak highlight
      if (r.streaks && r.streaks.length > 0) {
        const best = r.streaks.reduce((a, b) => (a.current_streak || 0) >= (b.current_streak || 0) ? a : b, r.streaks[0]);
        html += `<div class="bg-gradient-to-br from-primary to-primary/60 rounded-xl p-6 text-slate-900 shadow-lg shadow-primary/20">
          <div class="flex flex-col items-center text-center">${Icons.flame(48, '#0b1215')}<h4 class="text-lg font-medium opacity-90 mt-2">Streak Terpanjang</h4>
          <p class="text-4xl font-black mt-1">${best.current_streak || 0} hari</p>
          <p class="text-sm mt-4 bg-black/10 px-4 py-2 rounded-xl backdrop-blur-md border border-black/5">Habit: ${escapeHtml(best.habit_name || '')}</p></div></div>`;
      }

      // Insights with per-type icons, deduplicated against Area Perbaikan
      const weakHabits = habitScores.filter(h => h.score < 70).sort((a, b) => a.score - b.score);
      const weakHabitIds = new Set(weakHabits.map(h => h.habit_id));
      if (r.insights && r.insights.length > 0) {
        const filteredInsights = r.insights.filter(i => !(i.type === 'declining' && i.habit_id && weakHabitIds.has(i.habit_id)));
        if (filteredInsights.length > 0) {
          const insightIcons = { best_day: Icons.trophy(20, '#fbbf24'), declining: Icons.trendingDown(20, '#f87171'), consistency: Icons.star(20, '#2ec2b3'), encouragement: Icons.rocket(20, '#a78bfa') };
          html += `<div class="bg-slate-900/50 rounded-xl p-6 border border-slate-800 relative overflow-hidden">
            <div class="absolute top-0 right-0 p-2 opacity-10">${Icons.lightbulb(60)}</div>
            <div class="flex items-start gap-4"><div class="bg-amber-900/40 p-3 rounded-lg text-amber-400">${Icons.lightbulb(24, '#fbbf24')}</div>
            <div class="flex-1"><h4 class="font-bold text-white mb-1">Insight Minggu Ini</h4>`;
          filteredInsights.forEach(i => {
            const icon = insightIcons[i.type] || '';
            html += `<div class="flex items-start gap-2 mt-3"><span class="shrink-0 mt-0.5">${icon}</span><p class="text-sm text-slate-400 leading-relaxed">${escapeHtml(i.message)}</p></div>`;
          });
          html += '</div></div></div>';
        }
      }

      // Area Perbaikan - all habits below 70% with actionable tips
      if (weakHabits.length > 0) {
        html += `<div class="bg-slate-900/50 rounded-xl p-6 border-l-4 border-red-500">
          <div class="flex items-start gap-4"><div class="bg-red-900/40 p-3 rounded-lg text-red-400">${Icons.alertTriangle(24, '#f87171')}</div>
          <div class="flex-1"><h4 class="font-bold text-white mb-1">Area Perbaikan</h4><p class="text-sm text-slate-400 mb-3">Habit yang perlu perhatian lebih:</p><div class="space-y-2">`;
        weakHabits.forEach(h => {
          const tip = h.score === 0 ? 'Mulai dengan langkah kecil hari ini!' : h.score < 25 ? 'Coba pasang pengingat untuk membantu konsistensi.' : h.score < 50 ? 'Sudah ada progres! Tingkatkan frekuensinya.' : 'Hampir konsisten! Sedikit lagi untuk jadi kebiasaan.';
          html += `<div class="p-3 bg-slate-800 rounded-lg border border-slate-700"><div class="flex items-center justify-between"><span class="font-bold text-white">${escapeHtml(h.habit_name)}</span><span class="text-red-500 font-bold text-sm">${Math.round(h.score)}%</span></div><p class="text-xs text-slate-500 mt-1">${tip}</p></div>`;
        });
        html += '</div></div></div></div>';
      }

      html += '</div></div>';

      // Weekly grid — daily breakdown from real data
      const bd = r.daily_breakdown || [];
      if (bd.length > 0) {
        html += `<div class="mt-8 bg-slate-900/50 rounded-xl p-6 border border-slate-800"><h3 class="text-lg font-bold text-white mb-6">Ringkasan Harian</h3><div class="space-y-3">`;
        bd.forEach(day => {
          const pct = Math.round(day.rate);
          const barColor = pct >= 100 ? 'bg-emerald-500' : pct >= 50 ? 'bg-primary' : pct > 0 ? 'bg-amber-500' : 'bg-slate-700';
          const label = day.day_name + ' (' + formatDateShort(day.date) + ')';
          html += `<div class="space-y-1.5">
            <div class="flex justify-between items-center text-sm">
              <span class="font-medium text-slate-300">${escapeHtml(label)}</span>
              <span class="text-xs text-slate-400">${day.completed}/${day.total} <span class="font-bold ${pct >= 100 ? 'text-emerald-400' : 'text-slate-300'}">${pct}%</span></span>
            </div>
            <div class="h-2 bg-slate-800 rounded-full overflow-hidden"><div class="h-full ${barColor} rounded-full progress-fill" style="width:${pct}%"></div></div>
          </div>`;
        });
        html += '</div></div>';
      }
      html += `<footer class="py-10 text-center text-slate-400 text-sm"><p>&copy; 2026 HabitFlow Tracker. Dikembangkan untuk progresmu.</p></footer>`;

      content.innerHTML = html;
    } catch (err) {
      content.innerHTML = `<div class="text-center py-12">${Icons.emptyReport()}<p class="text-slate-500 text-sm mt-4">${escapeHtml(err.message) || 'Belum ada data laporan'}</p></div>`;
    }
  }

  // ════════════════════════════════════════════════════════════════
  // PAGE: SETTINGS
  // ════════════════════════════════════════════════════════════════
  async function renderSettings() {
    const user = API.getUser();
    const pushSupported = 'PushManager' in window;
    let pushEnabled = false;
    if (pushSupported && navigator.serviceWorker?.controller) {
      try { const reg = await navigator.serviceWorker.ready; const sub = await reg.pushManager.getSubscription(); pushEnabled = !!sub; } catch {}
    }
    const notificationStatusText = !pushSupported
      ? 'Tidak Didukung'
      : (pushEnabled ? 'Aktif' : 'Nonaktif');
    const notificationStatusClass = !pushSupported
      ? 'bg-amber-500/15 text-amber-300 border border-amber-400/30'
      : (pushEnabled
        ? 'bg-emerald-500/15 text-emerald-300 border border-emerald-400/30'
        : 'bg-slate-700/70 text-slate-300 border border-slate-600/80');

    const headerHTML = `
      <div>
        <h1 class="text-xl font-bold text-white">Pengaturan</h1>
        <p class="text-xs text-slate-400 mt-1">Kelola akun, notifikasi, dan preferensi aplikasi</p>
      </div>
      <div></div>`;
    const contentHTML = `
      <div class="max-w-2xl space-y-4">
        <div class="hf-settings-panel rounded-2xl p-5 md:p-6">
          <div class="flex items-center gap-4">
            <div class="size-14 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-xl">${(user?.name || 'U')[0].toUpperCase()}</div>
            <div class="min-w-0 flex-1">
              <p class="font-bold text-white text-base truncate">${escapeHtml(user?.name || 'User')}</p>
              <p class="text-sm text-slate-400 truncate">${escapeHtml(user?.email || '')}</p>
            </div>
            <span class="text-[11px] px-2.5 py-1 rounded-full bg-primary/15 text-primary border border-primary/30">Akun Aktif</span>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="hf-settings-card rounded-2xl p-5">
            <div class="flex items-start justify-between gap-3 mb-4">
              <div>
                <h3 class="font-bold text-sm text-white">Notifikasi</h3>
                <p class="text-xs text-slate-400 mt-1">Pengingat habit harian</p>
              </div>
              <span class="text-[11px] px-2.5 py-1 rounded-full ${notificationStatusClass}">${notificationStatusText}</span>
            </div>
          ${pushSupported ? `
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm text-slate-100 font-medium">Push Notification</p>
                <p class="text-xs text-slate-500 mt-0.5">Aktifkan agar tidak lupa check-in</p>
              </div>
              <button id="btn-push-toggle" role="switch" aria-checked="${pushEnabled}" aria-label="Toggle push notification" class="hf-toggle tap-target ${pushEnabled ? 'is-on' : ''}">
                <span class="hf-toggle-knob"></span>
              </button>
            </div>` : '<p class="text-sm text-slate-500">Browser tidak mendukung push notification</p>'}
          </div>

          <div class="hf-settings-card rounded-2xl p-5">
            <h3 class="font-bold text-sm text-white mb-3">Informasi Aplikasi</h3>
            <div class="space-y-2.5 text-xs">
              <div class="hf-settings-row"><span class="text-slate-400">Versi</span><span class="text-slate-200 font-medium">HabitFlow v1.0</span></div>
              <div class="hf-settings-row"><span class="text-slate-400">Platform</span><span class="text-slate-200 font-medium">Web Progressive App</span></div>
              <div class="hf-settings-row"><span class="text-slate-400">Proyek</span><span class="text-slate-200 font-medium">Skripsi S1 TI, UBSI</span></div>
            </div>
          </div>
        </div>

        <div class="hf-settings-card rounded-2xl p-5">
          <h3 class="font-bold text-sm text-white mb-2">Sesi</h3>
          <p class="text-xs text-slate-400 mb-4">Keluar dari akun ini jika kamu selesai menggunakan perangkat.</p>
          <button id="btn-logout" class="w-full py-3 bg-slate-900/70 border border-red-500/35 text-red-400 font-semibold rounded-xl hover:bg-red-900/20 transition tap-target">
            <span class="flex items-center justify-center gap-2">${Icons.logOut(20, '#f87171')} Keluar</span>
          </button>
        </div>
      </div>`;
    app().innerHTML = appShell('settings', headerHTML, contentHTML, { fab: false });
    $('#btn-logout').onclick = async () => {
      if (!confirm('Yakin ingin keluar?')) return;
      await API.logout();
      toast('Berhasil keluar', 'info');
    };
    if (pushSupported && $('#btn-push-toggle')) {
      $('#btn-push-toggle').onclick = async function () {
        try {
          const reg = await navigator.serviceWorker.ready;
          let sub = await reg.pushManager.getSubscription();
          if (sub) { await API.unsubscribePush(sub.endpoint); await sub.unsubscribe(); toast('Notifikasi dimatikan', 'info'); }
          else {
            if ('Notification' in window && Notification.permission === 'default') {
              const permission = await Notification.requestPermission();
              if (permission !== 'granted') throw new Error('Izin notifikasi ditolak');
            } else if ('Notification' in window && Notification.permission === 'denied') {
              throw new Error('Izin notifikasi diblokir di browser');
            }

            const vapidRes = await API.getVAPIDKey(); const key = vapidRes.data.vapid_public_key;
            const normalizedKey = normalizeVAPIDPublicKey(key);
            if (!normalizedKey) throw new Error('VAPID public key tidak valid. Cek konfigurasi server (.env).');

            sub = await reg.pushManager.subscribe({ userVisibleOnly: true, applicationServerKey: urlBase64ToUint8Array(normalizedKey) });
            const subJson = sub.toJSON();
            await API.subscribePush({ endpoint: sub.endpoint, keys: { p256dh: subJson.keys.p256dh, auth: subJson.keys.auth } });
            toast('Notifikasi diaktifkan!');
          }
          renderSettings();
        } catch (err) { toast('Gagal: ' + err.message, 'error'); }
      };
    }
  }

  function normalizeVAPIDPublicKey(value) {
    if (typeof value !== 'string') return '';
    const key = value.trim().replace(/^['"]|['"]$/g, '');
    if (!key) return '';
    if (key.toLowerCase().includes('your-vapid')) return '';
    if (!/^[A-Za-z0-9_-]+$/.test(key)) return '';
    if (key.length < 60) return '';
    return key;
  }

  function urlBase64ToUint8Array(base64String) {
    if (!base64String || typeof base64String !== 'string') {
      throw new Error('VAPID key kosong atau bukan string');
    }

    if (base64String.length % 4 === 1) {
      throw new Error('Panjang VAPID key tidak valid untuk Base64');
    }

    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');

    let rawData = '';
    try {
      rawData = window.atob(base64);
    } catch {
      throw new Error('VAPID key tidak dapat di-decode (format base64url tidak valid)');
    }

    const outputArray = new Uint8Array(rawData.length);
    for (let i = 0; i < rawData.length; ++i) outputArray[i] = rawData.charCodeAt(i);
    return outputArray;
  }

})();
