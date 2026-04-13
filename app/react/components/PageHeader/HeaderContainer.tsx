import { PropsWithChildren, createContext, useContext } from 'react';

const Context = createContext<null | boolean>(null);
Context.displayName = 'PageHeaderContext';

export function useHeaderContext() {
  const context = useContext(Context);

  if (context == null) {
    throw new Error('Should be nested inside a HeaderContainer component');
  }
}
interface Props {
  id?: string;
}

export function HeaderContainer({ id, children }: PropsWithChildren<Props>) {
  return (
    <Context.Provider value>
      <div
        id={id}
        className="row min-h-[60px] bg-[var(--bg-widget-color)] !mb-[5px] !rounded-none !border-0 !border-b !border-solid !border-b-[var(--border-widget)] !shadow-none"
      >
        <div id="loadingbar-placeholder" />
        <div className="col-xs-12">
          <div className="flex items-center justify-between [&_div]:truncate">
            {children}
          </div>
        </div>
      </div>
    </Context.Provider>
  );
}
