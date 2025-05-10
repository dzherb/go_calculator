export const timeout = (ms) => new Promise((resolve) => setTimeout(() => resolve(),  ms))

export const snakeToCamel = str =>
  str.toLowerCase().replace(/([-_][a-z])/g, group =>
    group
      .toUpperCase()
      .replace('-', '')
      .replace('_', '')
  );
