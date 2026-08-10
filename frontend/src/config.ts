// Build-time configuration, injected by webpack's DefinePlugin
// (see webpack/base.config.js) from frontend/.env.
//
// These values are inlined into the browser bundle and are therefore public —
// keeping them out of the source tree makes them easy to rotate per environment,
// it does not make them secret. Access has to be restricted at the provider:
// HTTP referrer limits for the Maps key, allowed origins for EmailJS.

export const GOOGLE_MAPS_API_KEY = process.env.GOOGLE_MAPS_API_KEY || '';

export const EMAILJS_SERVICE_ID = process.env.EMAILJS_SERVICE_ID || '';
export const EMAILJS_TEMPLATE_ID = process.env.EMAILJS_TEMPLATE_ID || '';
export const EMAILJS_PUBLIC_KEY = process.env.EMAILJS_PUBLIC_KEY || '';

export const GA_TRACKING_ID = process.env.GA_TRACKING_ID || '';
