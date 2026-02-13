/**
 * Biến API (giống Postman): Base URL và Access token.
 * Bấm nút **Copy** có sẵn phía trên khối cURL → copy ra đã thay Base URL & token, dùng được ngay.
 */
(function () {
  var BASE_PLACEHOLDER = 'https://api.example.com';
  var TOKEN_PLACEHOLDER = '<access_token>';
  var KEY_BASE = 'api_docs_base_url';
  var KEY_TOKEN = 'api_docs_access_token';

  function getBase() {
    var v = localStorage.getItem(KEY_BASE);
    return (v && v.trim()) ? v.trim() : BASE_PLACEHOLDER;
  }

  function getToken() {
    var v = localStorage.getItem(KEY_TOKEN);
    return (v && v.trim()) ? v.trim() : TOKEN_PLACEHOLDER;
  }

  function replaceInText(text) {
    return (text || '')
      .split(BASE_PLACEHOLDER).join(getBase())
      .split(TOKEN_PLACEHOLDER).join(getToken());
  }

  function showToast(msg) {
    var prev = document.getElementById('api-toast');
    if (prev) prev.remove();
    var el = document.createElement('div');
    el.id = 'api-toast';
    el.className = 'api-toast';
    el.textContent = msg;
    document.body.appendChild(el);
    setTimeout(function () {
      el.classList.add('api-toast--out');
      setTimeout(function () { el.remove(); }, 220);
    }, 1500);
  }

  /** Cập nhật nội dung khối cURL trên UI theo Base URL & token hiện tại. */
  function replaceInDOM() {
    var pres = document.querySelectorAll('main pre');
    for (var i = 0; i < pres.length; i++) {
      var code = pres[i].querySelector('code') || pres[i];
      var text = code.textContent || '';
      var orig = code.getAttribute('data-api-original');
      if (orig) {
        code.textContent = replaceInText(orig);
      } else if (text.indexOf('curl') !== -1 && (text.indexOf(BASE_PLACEHOLDER) !== -1 || text.indexOf(TOKEN_PLACEHOLDER) !== -1)) {
        code.setAttribute('data-api-original', text);
        code.textContent = replaceInText(text);
      }
    }
  }

  function onCopyClick(e) {
    var button = e.target && e.target.closest && e.target.closest('button');
    if (!button) return;

    var block = button.closest('.theme-code-block') || button.closest('[class*="codeBlock"]') || button.closest('[class*="CodeBlock"]');
    if (!block) return;

    var pre = block.querySelector('pre');
    if (!pre) return;

    var code = pre.querySelector('code') || pre;
    var text = code.textContent || '';
    var orig = code.getAttribute('data-api-original') || '';
    if (text.indexOf('curl') === -1 && orig.indexOf('curl') === -1) return;
    if (!orig && text.indexOf(BASE_PLACEHOLDER) === -1 && text.indexOf(TOKEN_PLACEHOLDER) === -1) return;

    var aria = (button.getAttribute('aria-label') || '').toLowerCase();
    var title = (button.getAttribute('title') || '').toLowerCase();
    var isCopy = /copy|clipboard|sao chép/.test(aria) || /copy|clipboard|sao chép/.test(title);
    if (!isCopy) {
      var btns = block.querySelectorAll('button');
      isCopy = btns.length > 0 && button === btns[0];
    }
    if (!isCopy) return;

    e.preventDefault();
    e.stopImmediatePropagation();

    var out = replaceInText(orig || text);
    navigator.clipboard.writeText(out).then(
      function () {
        showToast('Đã sao chép!');
        var lbl = button.getAttribute('aria-label');
        if (lbl) button.setAttribute('aria-label', 'Đã copy!');
        setTimeout(function () { if (lbl) button.setAttribute('aria-label', lbl); }, 1500);
      }
    );
  }

  function save() {
    var baseEl = document.getElementById('api-base-url');
    var tokenEl = document.getElementById('api-access-token');
    if (baseEl) localStorage.setItem(KEY_BASE, baseEl.value.trim());
    if (tokenEl) localStorage.setItem(KEY_TOKEN, tokenEl.value);
    replaceInDOM();
    var status = document.getElementById('api-config-status');
    if (status) {
      status.textContent = 'Đã lưu. URL & token đã cập nhật trên toàn bộ khối cURL.';
      setTimeout(function () { status.textContent = ''; }, 3000);
    }
  }

  function init() {
    var baseEl = document.getElementById('api-base-url');
    var tokenEl = document.getElementById('api-access-token');
    var btn = document.getElementById('api-save-config');

    if (baseEl) baseEl.value = localStorage.getItem(KEY_BASE) || '';
    if (tokenEl) tokenEl.value = localStorage.getItem(KEY_TOKEN) || '';
    if (btn) btn.addEventListener('click', save);

    document.addEventListener('click', onCopyClick, true);

    setTimeout(replaceInDOM, 500);
    var main = document.querySelector('main');
    if (main) {
      var mo = new MutationObserver(function () { setTimeout(replaceInDOM, 150); });
      mo.observe(main, { childList: true, subtree: true });
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
