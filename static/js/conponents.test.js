import {FullPathBuilder, Header, PathsHolder, ParameterSymbols} from './components.js';

const consoleSpy = jest.spyOn(console, 'log').mockImplementation();

function mockSymbolsGetter(data) {
    return (callback) => {callback(data)};
}

const goodSymbolsGetter = mockSymbolsGetter({"Symbols":[{"Parameter":"offset","Symbol":"σ"},{"Parameter":"phi","Symbol":"φ"},{"Parameter":"frequency","Symbol":"ω"}]});

describe('Unit tests for looking up symbols for parameter names', () => {
    beforeEach(() => {
      consoleSpy.mockClear()
    })

    test('Lookup fails gracefully when initialised with badly formatted data', () => {
        const badSymbolsGetter = mockSymbolsGetter({});

        const symbols = new ParameterSymbols(badSymbolsGetter);

        expect(console.log).toBeCalledTimes(1);
    })

    test('Lookup returns correct symbol when available', () => {
        const symbols = new ParameterSymbols(goodSymbolsGetter);

        expect(symbols.lookup('phi')).toBe('φ');

        expect(console.log).toBeCalledTimes(0);
    })

    // test('Lookup returns parameter when no symbol available', () => {
    //     const symbols = new ParameterSymbols(goodSymbolsGetter);

    //     const parameter = 'Beta';

    //     expect(symbols.lookup(parameter)).toBe(parameter);

    //     expect(console.log).toBeCalledTimes(0);
    // })
})