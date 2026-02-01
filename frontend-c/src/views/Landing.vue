<template>
  <div class="landing">
    <section class="hero">
      <div class="hero-bg" aria-hidden="true" />
      <div class="hero-content">
        <h1 class="hero-title">VibeCoding<span class="hero-title-accent">社区</span></h1>
        <p class="hero-tagline">集内容、社交、娱乐、技术于一体 · 多源化创作与交流空间</p>
        <div class="hero-actions">
          <router-link to="/posts" class="btn btn-primary btn-hero">
            <Icons name="feed" size="sm" /> 浏览帖子
          </router-link>
          <router-link v-if="!token" to="/login" class="btn btn-ghost btn-hero">登录</router-link>
          <router-link v-else to="/publish" class="btn btn-ghost btn-hero">
            <Icons name="pen" size="sm" /> 发帖
          </router-link>
        </div>
      </div>
    </section>

    <section class="features">
      <h2 class="features-heading">在这里，你可以</h2>
      <div class="features-grid">
        <article class="feature-card" v-for="(f, i) in features" :key="f.title">
          <div class="feature-icon" :style="{ animationDelay: `${i * 0.1}s` }">
            <Icons :name="f.icon" size="lg" />
          </div>
          <h3 class="feature-title">{{ f.title }}</h3>
          <p class="feature-desc">{{ f.desc }}</p>
        </article>
      </div>
    </section>

    <section class="cta">
      <p class="cta-text">发现灵感 · 分享生活 · 连接同好</p>
      <router-link to="/posts" class="btn btn-primary btn-lg">立即探索帖子</router-link>
    </section>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import Icons from '../components/Icons.vue'

const auth = useAuthStore()
const token = computed(() => auth.token)

const features = [
  { icon: 'pen', title: '内容创作', desc: '图文、视频随心发，记录与分享你的想法与作品。' },
  { icon: 'chat', title: '社交互动', desc: '点赞、评论、消息提醒，与创作者和同好实时交流。' },
  { icon: 'heart', title: '娱乐发现', desc: '浏览热门与推荐，在轻松氛围里发现有趣内容。' },
  { icon: 'search', title: '技术交流', desc: '话题与搜索直达，和开发者、创作者一起成长。' },
]
</script>

<style scoped>
.landing {
  min-height: 100vh;
  width: 100%;
  min-width: 0;
}
/* Hero */
.hero {
  position: relative;
  min-height: 52vh;
  min-height: 52dvh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 24px 64px;
  padding-left: max(24px, env(safe-area-inset-left));
  padding-right: max(24px, env(safe-area-inset-right));
  overflow: hidden;
}
@media (min-width: 769px) {
  .hero {
    min-height: 72vh;
    min-height: 72dvh;
    padding: 80px 48px 100px;
    padding-left: max(48px, env(safe-area-inset-left));
    padding-right: max(48px, env(safe-area-inset-right));
  }
}
.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, var(--color-primary-light) 0%, rgba(129, 216, 207, 0.08) 50%, var(--color-primary-light) 100%);
  background-size: 200% 200%;
  animation: gradientShift 12s ease infinite;
}
@keyframes gradientShift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}
.hero-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 560px;
}
@media (min-width: 769px) {
  .hero-content { max-width: 680px; }
}
.hero-title {
  font-family: var(--font-heading);
  font-size: clamp(2rem, 6vw, 3rem);
  font-weight: 700;
  margin: 0 0 16px;
  color: var(--color-text);
  letter-spacing: -0.02em;
}
@media (min-width: 769px) {
  .hero-title { font-size: clamp(2.5rem, 4vw, 3.5rem); margin-bottom: 20px; }
}
.hero-title-accent {
  color: var(--color-primary);
  font-weight: 800;
}
.hero-tagline {
  font-size: clamp(0.875rem, 2.5vw, 1rem);
  color: var(--color-text-secondary);
  margin: 0 0 32px;
  line-height: 1.6;
}
@media (min-width: 769px) {
  .hero-tagline { font-size: 1.125rem; margin-bottom: 40px; }
}
.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}
@media (min-width: 769px) {
  .hero-actions { gap: 16px; }
}
.btn-hero { padding: 12px 24px; font-size: 1rem; }
@media (min-width: 769px) {
  .btn-hero { padding: 14px 28px; font-size: 1.0625rem; }
}
.btn-hero .icon-wrap { margin-right: 4px; }

/* Features */
.features {
  padding: 56px 24px 64px;
  padding-left: max(24px, env(safe-area-inset-left));
  padding-right: max(24px, env(safe-area-inset-right));
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
  min-width: 0;
}
@media (min-width: 769px) {
  .features {
    max-width: 1200px;
    padding: 80px 48px 100px;
    padding-left: max(48px, env(safe-area-inset-left));
    padding-right: max(48px, env(safe-area-inset-right));
  }
}
.features-heading {
  font-family: var(--font-heading);
  font-size: clamp(1.25rem, 4vw, 1.5rem);
  font-weight: 600;
  text-align: center;
  margin: 0 0 40px;
  color: var(--color-text);
}
@media (min-width: 769px) {
  .features-heading { font-size: 1.75rem; margin-bottom: 56px; }
}
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
}
@media (min-width: 769px) {
  .features-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 32px;
  }
}
.feature-card {
  background: var(--color-card);
  border-radius: var(--radius);
  padding: 28px 20px;
  text-align: center;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow);
  transition: box-shadow var(--transition), border-color var(--transition), transform var(--transition);
}
@media (min-width: 769px) {
  .feature-card { padding: 36px 24px; }
}
.feature-card:hover {
  box-shadow: var(--shadow-hover);
  border-color: rgba(129, 216, 207, 0.35);
  transform: translateY(-2px);
}
.feature-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-light);
  color: var(--color-primary-dark);
  border-radius: 50%;
  animation: fadeInUp 0.6s ease backwards;
}
@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
.feature-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0 0 8px;
  color: var(--color-text);
}
.feature-desc {
  font-size: 14px;
  color: var(--color-text-muted);
  margin: 0;
  line-height: 1.5;
}

/* CTA */
.cta {
  padding: 48px 24px 64px;
  padding-left: max(24px, env(safe-area-inset-left));
  padding-right: max(24px, env(safe-area-inset-right));
  text-align: center;
  background: linear-gradient(180deg, transparent 0%, var(--color-primary-light) 100%);
  border-radius: var(--radius) var(--radius) 0 0;
}
@media (min-width: 769px) {
  .cta {
    padding: 72px 48px 96px;
    padding-left: max(48px, env(safe-area-inset-left));
    padding-right: max(48px, env(safe-area-inset-right));
  }
}
.cta-text {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0 0 20px;
}
@media (min-width: 769px) {
  .cta-text { font-size: 1.125rem; margin-bottom: 28px; }
}
.btn-lg { padding: 14px 32px; font-size: 1.05rem; }
@media (min-width: 769px) {
  .btn-lg { padding: 16px 40px; font-size: 1.125rem; }
}

@media (max-width: 768px) {
  .hero { padding: 32px 16px 48px; padding-left: max(16px, env(safe-area-inset-left)); padding-right: max(16px, env(safe-area-inset-right)); }
  .features { padding: 40px 16px 48px; padding-left: max(16px, env(safe-area-inset-left)); padding-right: max(16px, env(safe-area-inset-right)); }
  .features-grid { grid-template-columns: 1fr; gap: 16px; }
  .feature-card { padding: 24px 16px; }
  .cta { padding: 40px 16px 56px; padding-left: max(16px, env(safe-area-inset-left)); padding-right: max(16px, env(safe-area-inset-right)); }
}
@media (max-width: 375px) {
  .hero { padding: 24px 12px 40px; }
  .hero-actions { gap: 8px; }
  .btn-hero { padding: 10px 18px; font-size: 0.9375rem; }
  .features { padding: 32px 12px 40px; }
  .cta { padding: 32px 12px 48px; }
}
@media (prefers-reduced-motion: reduce) {
  .feature-card:hover { transform: none; }
}
</style>
