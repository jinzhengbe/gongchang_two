<template>
  <footer class="footer">
    <div class="footer-content">
      <div class="footer-section">
        <h3>{{ $t('footer.about') }}</h3>
        <p>{{ $t('footer.description') }}</p>
      </div>
      
      <div class="footer-section">
        <h3>{{ $t('footer.contact') }}</h3>
        <p>
          {{ $t('footer.email') }}: support@sewingmast.com<br>
          {{ $t('footer.phone') }}: +86-XXX-XXXX-XXXX<br>
          {{ $t('footer.website') }}: https://sewingmast.com
        </p>
      </div>

      <div class="footer-section">
        <h3>{{ $t('footer.links') }}</h3>
        <ul>
          <li><a href="/about">{{ $t('footer.aboutUs') }}</a></li>
          <li><a href="/terms">{{ $t('footer.terms') }}</a></li>
          <li><a href="/privacy">{{ $t('footer.privacy') }}</a></li>
        </ul>
      </div>
    </div>

    <div class="footer-bottom">
      <p>© {{ currentYear }} SewingMast. {{ $t('footer.copyright') }}</p>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>

<style scoped>
.footer {
  background: #f8f9fa;
  padding: 3rem 1rem 1rem;
  margin-top: auto;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.footer-section h3 {
  margin-bottom: 1rem;
  color: #333;
}

.footer-section ul {
  list-style: none;
  padding: 0;
}

.footer-section ul li {
  margin-bottom: 0.5rem;
}

.footer-section a {
  color: #666;
  text-decoration: none;
}

.footer-section a:hover {
  color: #333;
}

.footer-bottom {
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #ddd;
  text-align: center;
  color: #666;
}

/* Responsive styles */
@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .footer {
    padding: 2rem 1rem 1rem;
  }
}
</style> 