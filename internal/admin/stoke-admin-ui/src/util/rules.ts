
export const require = function(msg : string) : (v : string) => string | boolean {
  return (v : string) => !!v || `${msg} is required.`
}
