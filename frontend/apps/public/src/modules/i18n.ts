import i18n from "i18next";
import { initReactI18next } from "react-i18next";

const resources = {
  en: {
    translation: {
      welcome: "Welcome to the personal website",
      intro: "This SPA scaffold is ready for customization."
    }
  },
  ja: {
    translation: {
      welcome: "個人サイトへようこそ",
      intro: "この SPA スケルトンはカスタマイズ可能です。"
    }
  }
};

void i18n.use(initReactI18next).init({
  resources,
  fallbackLng: "en",
  interpolation: {
    escapeValue: false
  }
});

export default i18n;
