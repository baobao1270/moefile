import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import en from '@/locales/en.json';
import zh_CN from '@/locales/zh_CN.json';


i18n.use(LanguageDetector).use(initReactI18next).init({
  fallbackLng: 'en',
  debug: import.meta.env.NODE_ENV === 'development',
  interpolation: { escapeValue: false },
  resources: {
    'en': { ...en },
    'zh-CN': { ...zh_CN },
  },
})
