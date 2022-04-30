const palette = {
  white: '#FFFFFF',
  green: '#388E3C',
  grey: '#101010',
  grey2: '#F5F5F5',
  orange: '#FB8C00'
};

export const colors = {
  background: {
    white: palette.white,
    green: palette.green,
    grey: palette.grey2,
  },
  text: {
    default: palette.grey,
    white: palette.white,
  },
  button: {
    default: palette.green,
    carousel: palette.orange,
  }
};

declare global {
  type ThemeColors =
    | 'background.white'
    | 'background.green'
    | 'background.grey2'
    | 'text.default'
    | 'text.white'
    | 'button.default'
    | 'button.carousel';
}
