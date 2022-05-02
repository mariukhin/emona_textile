const palette = {
  white: '#FFFFFF',
  green: '#388E3C',
  grey: '#101010',
  grey2: '#F5F5F5',
  orange: '#FB8C00'
};

export const colors = {
  background: {
    default: palette.grey2,
    white: palette.white,
    green: palette.green,
  },
  text: {
    default: palette.grey,
    white: palette.white,
    green: palette.green,
  },
  button: {
    default: palette.green,
    carousel: palette.orange,
  }
};

declare global {
  type ThemeColors =
    | 'background.default'
    | 'background.white'
    | 'background.green'
    | 'text.default'
    | 'text.white'
    | 'text.green'
    | 'button.default'
    | 'button.carousel';
}
